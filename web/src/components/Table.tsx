import { type JSX, For, Show } from 'solid-js'
import { cn } from '../lib/utils'

export type Column<T> = {
  key: string
  label: string
  render?: (item: T) => JSX.Element
  className?: string
}

type TableProps<T> = {
  columns: Column<T>[]
  data: T[]
  onRowClick?: (item: T) => void
  isLoading?: boolean
  emptyMessage?: string
}

export function Table<T>(props: TableProps<T>) {
  return (
    <div class='overflow-x-auto border border-gray-200 rounded-lg'>
      <table class='min-w-full divide-y divide-gray-200'>
        <thead class='bg-gray-50'>
          <tr>
            <For each={props.columns}>
              {(column) => (
                <th
                  scope='col'
                  class={cn(
                    'px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider',
                    column.className,
                  )}
                >
                  {column.label}
                </th>
              )}
            </For>
          </tr>
        </thead>
        <tbody class='bg-white divide-y divide-gray-200'>
          <Show
            when={!props.isLoading && props.data.length > 0}
            fallback={
              <tr>
                <td colspan={props.columns.length} class='px-4 py-8 text-center text-sm text-gray-500'>
                  {props.isLoading ? 'Loading...' : props.emptyMessage || 'No data available'}
                </td>
              </tr>
            }
          >
            <For each={props.data}>
              {(item) => (
                <tr
                  class={cn('hover:bg-gray-50 transition-colors', props.onRowClick && 'cursor-pointer')}
                  onClick={() => props.onRowClick?.(item)}
                >
                  <For each={props.columns}>
                    {(column) => (
                      <td class='px-4 py-3 whitespace-nowrap text-sm text-gray-900'>
                        {column.render ? column.render(item) : (item as any)[column.key]}
                      </td>
                    )}
                  </For>
                </tr>
              )}
            </For>
          </Show>
        </tbody>
      </table>
    </div>
  )
}
