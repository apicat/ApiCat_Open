import useTable from '@/hooks/useTable'
import { getIterationList, deleteIteration } from '@/api/iteration'
import { Iteration } from '@/typings'
import { AsyncMsgBox } from '@/components/AsyncMessageBox'
import { useI18n } from 'vue-i18n'

export type SelectedProjectKey = number | string | null

export const useIterationList = () => {
  const { t } = useI18n()
  const selectedProjectKeyRef = ref<SelectedProjectKey>(null)

  const queryParam: Record<string, any> = {
    project_id: selectedProjectKeyRef.value,
  }

  const { currentPage, getTableData, ...rest } = useTable<Iteration>(getIterationList, {
    searchParam: queryParam,
    isLoaded: false,
    dataKey: 'iterations',
  })

  const handleRemoveIteration = (iteration: Iteration) => {
    AsyncMsgBox({
      title: t('app.common.deleteTip'),
      content: '确定删除该迭代吗?',
      onOk: async () => {
        await deleteIteration({ iteration_public_id: iteration.id })
        await getTableData()
      },
    })
  }

  // 项目切换时获取当前项目的迭代列表
  watch(
    selectedProjectKeyRef,
    async () => {
      queryParam.project_id = selectedProjectKeyRef.value || ''
      currentPage.value = 1
      await getTableData()
    },
    {
      immediate: true,
    }
  )

  return {
    selectedProjectKeyRef,
    currentPage,
    ...rest,
    fetchIterationList: getTableData,
    handleRemoveIteration,
  }
}
