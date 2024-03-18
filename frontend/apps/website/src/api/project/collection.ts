import { delay, parseJSONWithDefault } from '@apicat/shared'
import type { AxiosRequestConfig } from 'axios'
import Ajax from '@/api/Ajax'
import { gatherSharedTokenWithParams } from '@/api/shareToken'
import { API_URL, queryStringify } from '@/commons'

// 获取集合列表
export function apiGetCollections(
  projectID: string,
  iterationID?: string,
): Promise<CollectionAPI.ResponseCollection[]> {
  const params: any = {}

  if (iterationID)
    params.iterationID = iterationID

  return Ajax.get(`/projects/${projectID}/collections`, { params: gatherSharedTokenWithParams(params, projectID) })
}

// 创建集合
export function apiCreateCollection(
  projectID: string,
  data: Omit<CollectionAPI.ResponseCollectionDetail, 'id'> & { iterationID?: string },
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return Ajax.post(`/projects/${projectID}/collections`, data)
}
export function apiAICreateCollection(
  projectID: string,
  data: ProjectAPI.RequestCreateCollectionWithAI,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return Ajax.post(`/projects/${projectID}/ai/collections`, data)
}

// 获取集合
export async function apiGetCollectDetail(
  projectID: string,
  collectionID: number,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  const collection: CollectionAPI.ResponseCollectionDetail = await Ajax.get(
    `/projects/${projectID}/collections/${collectionID}`,
    { params: gatherSharedTokenWithParams({}, projectID) },
    { isShowErrorMsg: false },
  )
  collection.content = parseJSONWithDefault(collection.content, [])
  return collection
}

// rename collection
export function apiRenameCollection(
  projectID: string,
  collection: CollectionAPI.ResponseCollection,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return apiEditCollectionDetail(projectID, collection as any)
}

// 编辑集合
export function apiEditCollectionDetail(
  projectID: string,
  collection: CollectionAPI.ResponseCollectionDetail,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  const { id, ...data } = collection
  return Ajax.put(`/projects/${projectID}/collections/${id}`, { ...data, content: JSON.stringify(data.content || '') })
}

// 删除集合
export function apiDeleteCollection(
  projectID: string,
  collectionID: number,
  data?: Record<string, any>,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return Ajax.delete(`/projects/${projectID}/collections/${collectionID}`, data)
}

// 集合排序
export function apiMoveCollection(
  projectID: string,
  data: CollectionAPI.RequestMoveCollection,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return Ajax.put(`/projects/${projectID}/collections/move`, data)
}

// 复制集合
export function apiCopyCollection(
  projectID: string,
  collectionID: number,
  data?: Record<string, any>,
): Promise<CollectionAPI.ResponseCollectionDetail> {
  return Ajax.post(`/projects/${projectID}/collections/${collectionID}/copy`, data)
}

// 导出集合
export function exportCollection({ projectID, collectionID, ...params }: any) {
  return `${API_URL}/project/${projectID}/collection/${collectionID}data${queryStringify(params)}`
}

export function apiExportCollection(
  projectID: string,
  collectionID: number,
  type: string,
  download: boolean,
): Promise<{ path: string }> {
  return Ajax.get(`/projects/${projectID}/collections/${collectionID}/export`, { params: { type, download } })
}

// 测试用例
export function apiCreateTestCase(
  projectID: string,
  collectionID: number | string,
  prompt?: string,
  regenerate?: boolean,
  config?: AxiosRequestConfig,
): Promise<void> {
  const data = prompt ? { prompt } : regenerate ? { regenerate } : undefined
  return Ajax.post(`/projects/${projectID}/collections/${collectionID}/testcases`, data, undefined, config)
}
export function apiRegenTestCaseList(projectID: string, collectionID: number | string, config?: AxiosRequestConfig) {
  return apiCreateTestCase(projectID, collectionID, undefined, true, config)
}

export function apiGetTestCaseList(projectID: string, collectionID: number | string, config?: AxiosRequestConfig): Promise<ProjectAPI.ResponseTestCase> {
  return Ajax.get(`/projects/${projectID}/collections/${collectionID}/testcases`, config)
}
export async function apiGetTestCaseDetail(
  projectID: string,
  collectionID: number | string,
  testCaseID: number,
): Promise<ProjectAPI.TestCaseDetail> {
  return Ajax.get(`/projects/${projectID}/collections/${collectionID}/testcases/${testCaseID}`)
}
export function apiReGenTestCase(
  projectID: string,
  collectionID: number,
  testCaseID: number,
  prompt: string,
): Promise<ProjectAPI.TestCaseDetail> {
  return Ajax.put(`/projects/${projectID}/collections/${collectionID}/testcases/${testCaseID}`, { prompt })
}
export function apiDeleteTestCase(projectID: string, collectionID: number, testCaseID: number): Promise<void> {
  return Ajax.delete(`/projects/${projectID}/collections/${collectionID}/testcases/${testCaseID}`)
}
