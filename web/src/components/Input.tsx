import { type JSX, splitProps } from 'solid-js'
import { cn } from '../lib/utils'

type InputProps = JSX.InputHTMLAttributes<HTMLInputElement> & {
  label?: string
  error?: string
}

export function Input(props: InputProps) {
  const [local, rest] = splitProps(props, ['label', 'error', 'class'])

  return (
    <div class='space-y-1'>
      {local.label && <label class='block text-sm font-medium text-gray-700'>{local.label}</label>}
      <input
        class={cn(
          'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
          local.error && 'border-red-500 focus:ring-red-500 focus:border-red-500',
          local.class,
        )}
        {...rest}
      />
      {local.error && <p class='text-sm text-red-600'>{local.error}</p>}
    </div>
  )
}
