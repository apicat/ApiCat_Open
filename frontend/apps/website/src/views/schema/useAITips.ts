import type { JSONSchemaTable } from '@apicat/components'
import axios from 'axios'
import { debounce } from 'lodash-es'
import { apiGetAIModel } from '@/api/project/definition/schema'

export function useAITips(
  project_id: string,
  schema: Ref<Definition.SchemaNode | null>,
  readonly: Ref<boolean>,
  updateSchema: (projectID: string, schema: Definition.SchemaNode) => Promise<void | undefined>,
  currentSchemaIDRef: Ref<string>,
) {
  const { escape, tab } = useMagicKeys()
  const jsonSchemaTableIns = ref<InstanceType<typeof JSONSchemaTable>>()

  const schemaName = ref('')
  const isAIMode = ref(false)
  const isShowAIStyle = ref(false)
  const isLoadingAICollection = ref(false)
  const preSchema = ref<Definition.SchemaNode | null>(null)

  // 避免请求后，文档不匹配问题
  const requestID = ref<string>()
  let abortController: AbortController | null = null

  // 不允许AI提示系列操作判断条件
  const notAllowAITips = () => preSchema.value?.id !== schema.value?.id || !schema.value || readonly.value

  // 获取AI提示数据
  async function getAITips() {
    if (notAllowAITips())
      return

    // 取消上次请求
    abortController?.abort()
    requestID.value = `${Date.now()},${schema.value!.id}`

    try {
      abortController = new AbortController()
      isLoadingAICollection.value = true
      const res = await apiGetAIModel(project_id, { requestID: unref(requestID), modelID: schema.value!.id, title: schema.value!.name }, { signal: abortController.signal })
      const { schema: aiSchema, requestID: resRequestID } = res || {}
      if (requestID.value === resRequestID && aiSchema) {
        schema.value!.schema = aiSchema
        isShowAIStyle.value = true
      }

      abortController = null
      // 重置请求标识
      isLoadingAICollection.value = false
      requestID.value = ''
    }
    catch (error: any) {
      // Cancelled Error 不需要重置
      if (!axios.isCancel(error)) {
        isLoadingAICollection.value = false
        requestID.value = ''
      }
    }
  }

  // 标题失去焦点时,延迟600避免title的debounce冲突
  const handleTitleBlur = debounce(() => cancelAITips(), 600)

  // 取消AI提示
  function cancelAITips() {
    // 取消上次请求
    abortController?.abort()

    // 切换了新的内容，不需要还原
    if (Number(currentSchemaIDRef.value) !== preSchema.value?.id)
      isAIMode.value = false

    // 还原文档
    if (isAIMode.value && preSchema.value && schema.value && Number(currentSchemaIDRef.value) === preSchema.value?.id) {
      schema.value.schema = { type: 'object' }
      isAIMode.value = false
    }

    // 重置
    requestID.value = ''
    isShowAIStyle.value = false
    isLoadingAICollection.value = false
  }

  // 确认AI提示
  function confirmAITips() {
    if (notAllowAITips())
      return
    // trigger watch name
    schema.value!.name = schema.value!.name
    isShowAIStyle.value = false
    isAIMode.value = false
    try {
      const newSchemaStr = JSON.stringify(schema.value)
      const copyPreSchemaStr = JSON.stringify(preSchema.value)
      if (newSchemaStr === copyPreSchemaStr)
        return

      updateSchema(project_id, schema.value!)
      // 更新preSchema
      preSchema.value = JSON.parse(newSchemaStr)
    }
    catch (e) {
      console.error('confirmAITips error', e)
    }
  }

  watch(schemaName, async () => {
    if (!schemaName.value)
      return

    if (isAIMode.value || jsonSchemaTableIns.value?.isEmpty()) {
      isAIMode.value = true
      // 保存之前确认的数据
      if (preSchema.value) {
        preSchema.value.name = schemaName.value
        preSchema.value.description = schema.value?.description
        await updateSchema(project_id, preSchema.value)
      }

      await getAITips()
    }
    else {
      isAIMode.value = false
    }
  })

  whenever(escape, () => {
    if (readonly.value || !isAIMode.value || !schema.value)
      return
    cancelAITips()
  })

  whenever(tab, () => {
    // tab触发时，获取AI数据，直接取消请求
    if (isLoadingAICollection.value)
      handleTitleBlur()

    if (readonly.value || !isAIMode.value || !schema.value)
      return

    confirmAITips()
  })

  // 点击文档区域,非编辑器内点击 -> 取消AI提示
  function onDocumentLayoutClick(e: MouseEvent) {
    // 允许点击的区域的dom path 路径含有.ac-schema-editor样式，有效点击
    if (isAIMode.value && !isLoadingAICollection.value && e.composedPath().find((el: any) => el.className?.includes('ac-schema-editor')))
      confirmAITips()
    else
      handleTitleBlur()
  }

  return {
    isAIMode,
    isShowAIStyle,
    jsonSchemaTableIns,
    schemaName,
    preSchema,

    handleTitleBlur,
    onDocumentLayoutClick,
    cancelAITips,
  }
}
