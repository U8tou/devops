package db

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"pkg/conf"

	"github.com/robfig/cron/v3"
)

var (
	backupCron   *cron.Cron
	backupCronMu sync.Mutex
)

// StartBackup 根据配置启动定时备份；enable 为 false 时不启动
func StartBackup() {
	cfg := conf.Db
	if !cfg.Backup.Enable {
		return
	}
	b := &cfg.Backup
	if b.Cron == "" {
		b.Cron = "0 0 0 * * *"
	}
	if b.Path == "" {
		b.Path = "./backup"
	}
	if b.KeepDays <= 0 {
		b.KeepDays = 7
	}
	if b.MaxCount <= 0 {
		b.MaxCount = 10
	}
	if b.MaxSize <= 0 {
		b.MaxSize = 1024
	}

	backupCronMu.Lock()
	if backupCron != nil {
		backupCron.Stop()
	}
	backupCron = cron.New(cron.WithSeconds())
	_, err := backupCron.AddFunc(b.Cron, func() {
		if err := DoBackup(); err != nil {
			log.Printf("[backup] error: %v", err)
		}
	})
	backupCronMu.Unlock()
	if err != nil {
		log.Printf("[backup] cron add func error: %v", err)
		return
	}
	backupCron.Start()
	log.Println("✅ db backup cron started:", b.Cron)
}

// StopBackup 停止定时备份
func StopBackup() {
	backupCronMu.Lock()
	if backupCron != nil {
		backupCron.Stop()
		backupCron = nil
	}
	backupCronMu.Unlock()
}

// DoBackup 执行一次完整备份（路径、清理、压缩、加密由配置决定）
func DoBackup() error {
	cfg := conf.Db
	b := &cfg.Backup
	dir, err := filepath.Abs(b.Path)
	if err != nil {
		return fmt.Errorf("backup path: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir backup: %w", err)
	}

	ext := ".sql"
	if b.Compress {
		ext = ".sql.gz"
	}
	if b.Encrypt {
		ext += ".enc"
	}
	name := fmt.Sprintf("%s_%s%s", cfg.UseDb, time.Now().Format("20060102_150405"), ext)
	outPath := filepath.Join(dir, name)

	var out io.Writer
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create backup file: %w", err)
	}
	defer f.Close()
	out = f

	if b.Compress {
		gw := gzip.NewWriter(out)
		defer gw.Close()
		out = gw
	}
	if b.Encrypt {
		if b.EncryptKey == "" {
			return fmt.Errorf("backup encrypt enabled but encryptKey is empty")
		}
		key := deriveKey(b.EncryptKey, b.EncryptLevel)
		encW, err := newEncryptWriter(out, key)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
		defer encW.Close()
		out = encW
	}

	eng := GetDb()
	if err = eng.DumpAll(out); err != nil {
		_ = os.Remove(outPath)
		return fmt.Errorf("xorm DumpAll: %w", err)
	}

	return cleanupBackups(dir, b.KeepDays, b.MaxCount, b.MaxSize)
}

// cleanupBackups 按 keepDays、maxCount、maxSize 清理旧备份
func cleanupBackups(dir string, keepDays, maxCount, maxSizeMB int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	maxBytes := int64(maxSizeMB) * 1024 * 1024
	cutoff := time.Now().AddDate(0, 0, -keepDays)

	type fileInfo struct {
		path string
		mod  time.Time
		size int64
	}
	var files []fileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(dir, e.Name())
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		files = append(files, fileInfo{path, info.ModTime(), info.Size()})
	}
	// 按修改时间倒序（新的在前）
	sort.Slice(files, func(i, j int) bool { return files[i].mod.After(files[j].mod) })

	var total int64
	keep := 0
	for _, f := range files {
		if f.mod.Before(cutoff) {
			_ = os.Remove(f.path)
			continue
		}
		if keep >= maxCount {
			_ = os.Remove(f.path)
			continue
		}
		if total+f.size > maxBytes {
			_ = os.Remove(f.path)
			continue
		}
		total += f.size
		keep++
	}
	return nil
}

func deriveKey(password, level string) []byte {
	key := sha256.Sum256([]byte(password))
	k := key[:]
	switch level {
	case "high":
		for i := 0; i < 10000; i++ {
			h := sha256.Sum256(k)
			k = h[:]
		}
	case "medium":
		for i := 0; i < 1000; i++ {
			h := sha256.Sum256(k)
			k = h[:]
		}
	}
	if len(k) > 32 {
		k = k[:32]
	}
	return k
}

func newEncryptWriter(w io.Writer, key []byte) (io.WriteCloser, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	if _, err := w.Write(iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	return &cipherWriter{stream: stream, w: w}, nil
}

type cipherWriter struct {
	stream cipher.Stream
	w     io.Writer
}

func (c *cipherWriter) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	c.stream.XORKeyStream(buf, p)
	return c.w.Write(buf)
}

func (c *cipherWriter) Close() error { return nil }
