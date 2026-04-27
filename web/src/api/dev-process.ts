import request from '@/utils/http'

export function fetchDevProcessList(params?: Api.DevProcess.ProcessSearchParams) {
  return request.get<Api.DevProcess.ProcessListResponse>({
    url: '/api/dev_process/list',
    params
  })
}

export function fetchDevProcessDetail(id: string) {
  return request.get<Api.DevProcess.ProcessDetail>({
    url: '/api/dev_process/get',
    params: { id }
  })
}

export function fetchDevProcessAdd(data: Api.DevProcess.ProcessAddReq) {
  return request.post<Api.DevProcess.ProcessAddRes>({
    url: '/api/dev_process/add',
    data
  })
}

export function fetchDevProcessEdit(data: Api.DevProcess.ProcessEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_process/edit',
    data
  })
}

export function fetchDevProcessEditFlow(data: Api.DevProcess.ProcessEditFlowReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_process/edit_flow',
    data
  })
}

export function fetchDevProcessEditEnv(data: Api.DevProcess.ProcessEditEnvReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_process/edit_env',
    data
  })
}

export function fetchDevProcessSetCronEnabled(data: Api.DevProcess.ProcessSetCronEnabledReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_process/set_cron_enabled',
    data
  })
}

/** 删除流程（后端软删：仅标记 delete_time） */
export function fetchDevProcessDel(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/dev_process/del',
    data: { ids }
  })
}

/** 录入时快速校验节点（Git 仓库 / SSH 等），与执行引擎共用 devflow 校验逻辑 */
export function fetchDevProcessValidateNode(data: Api.DevProcess.ValidateNodeReq) {
  return request.post<Api.DevProcess.ValidateNodeRes>({
    url: '/api/dev_process/validate_node',
    data,
    timeout: 35000,
    showErrorMessage: true
  })
}

/** 取消当前用户对指定流程的 SSE 流式执行 */
export function fetchDevProcessRunCancel(id: string) {
  return request.post<{ ok: boolean }>({
    url: `/api/dev_process/run_cancel?id=${encodeURIComponent(id)}`
  })
}

export function fetchDevProcessTagList() {
  return request.get<Api.DevProcess.TagListResponse>({
    url: '/api/dev_process/tag/list'
  })
}

export function fetchDevProcessTagAdd(data: Api.DevProcess.TagAddReq) {
  return request.post<Api.DevProcess.TagAddRes>({
    url: '/api/dev_process/tag/add',
    data
  })
}

export function fetchDevProcessTagEdit(data: Api.DevProcess.TagEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_process/tag/edit',
    data
  })
}

export function fetchDevProcessTagDel(id: string) {
  return request.del<{ affect: string }>({
    url: `/api/dev_process/tag/del?id=${encodeURIComponent(id)}`
  })
}
