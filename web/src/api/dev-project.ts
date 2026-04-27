import request from '@/utils/http'

export function fetchDevProjectList(params?: Api.DevProject.ProjectSearchParams) {
  return request.get<Api.DevProject.ProjectListResponse>({
    url: '/api/dev_project/list',
    params
  })
}

export function fetchDevProjectDetail(id: string) {
  return request.get<Api.DevProject.ProjectDetail>({
    url: '/api/dev_project/get',
    params: { id }
  })
}

export function fetchDevProjectAdd(data: Api.DevProject.ProjectAddReq) {
  return request.post<Api.DevProject.ProjectAddRes>({
    url: '/api/dev_project/add',
    data
  })
}

export function fetchDevProjectEdit(data: Api.DevProject.ProjectEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_project/edit',
    data
  })
}

export function fetchDevProjectEditMind(data: Api.DevProject.ProjectEditMindReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_project/edit_mind',
    data,
    timeout: 60000
  })
}

export function fetchDevProjectDel(ids: string[]) {
  return request.del<{ affect: string }>({
    url: '/api/dev_project/del',
    data: { ids }
  })
}

export function fetchDevProjectTagList() {
  return request.get<Api.DevProject.TagListResponse>({
    url: '/api/dev_project/tag/list'
  })
}

export function fetchDevProjectTagAdd(data: Api.DevProject.TagAddReq) {
  return request.post<Api.DevProject.TagAddRes>({
    url: '/api/dev_project/tag/add',
    data
  })
}

export function fetchDevProjectTagEdit(data: Api.DevProject.TagEditReq) {
  return request.put<{ affect: string }>({
    url: '/api/dev_project/tag/edit',
    data
  })
}

export function fetchDevProjectTagDel(id: string) {
  return request.del<{ affect: string }>({
    url: `/api/dev_project/tag/del?id=${encodeURIComponent(id)}`
  })
}
