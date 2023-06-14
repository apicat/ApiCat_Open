import type { RouteRecordRaw } from 'vue-router'
import { compile } from 'path-to-regexp'
import { PROJECT_DETAIL_PATH } from './constant'
import { MemberAuthorityInProject } from '@/typings/member'

export const DOCUMENT_DETAIL_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id?'
export const DOCUMENT_EDIT_PATH = PROJECT_DETAIL_PATH + '/doc/:doc_id/edit'

export const DOCUMENT_DETAIL_NAME = 'document.detail'
export const DOCUMENT_EDIT_NAME = 'document.edit'

export const getDocumentDetailPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id })
export const getDocumentEditPath = (project_id: number | string, doc_id: number | string) => compile(DOCUMENT_DETAIL_PATH)({ project_id, doc_id }) + '/edit'

const DocumentDetailPage = () => import('@/views/document/DocumentDetailPage.vue')
const DocumentEditPage = () => import('@/views/document/DocumentEditPage.vue')

export const documentDetailRoute: RouteRecordRaw = {
  name: DOCUMENT_DETAIL_NAME,
  path: DOCUMENT_DETAIL_PATH,
  component: DocumentDetailPage,
}

export const documentEditRoute: RouteRecordRaw = {
  name: DOCUMENT_EDIT_NAME,
  path: DOCUMENT_EDIT_PATH,
  component: DocumentEditPage,
  meta: {
    editableRoles: [MemberAuthorityInProject.MANAGER, MemberAuthorityInProject.WRITE],
  },
}
