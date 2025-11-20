import { type JSX, splitProps, For } from 'solid-js'
import { cn } from '../lib/utils'

type SelectOption = {
  value: string
  label: string
}

type SelectProps = JSX.SelectHTMLAttributes<HTMLSelectElement> & {
  label?: string
  error?: string
  options: SelectOption[]
  placeholder?: string
}

export function Select(props: SelectProps) {
  const [local, rest] = splitProps(props, ['label', 'error', 'class', 'options', 'placeholder'])

  return (
    <div class='space-y-1'>
      {local.label && <label class='block text-sm font-medium text-gray-700'>{local.label}</label>}
      <select
        class={cn(
          'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
          local.error && 'border-red-500 focus:ring-red-500 focus:border-red-500',
          local.class,
        )}
        {...rest}
      >
        {local.placeholder && <option value=''>{local.placeholder}</option>}
        <For each={local.options}>{(option) => <option value={option.value}>{option.label}</option>}</For>
      </select>
      {local.error && <p class='text-sm text-red-600'>{local.error}</p>}
    </div>
  )
}
