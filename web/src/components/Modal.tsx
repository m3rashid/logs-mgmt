import { createSignal, type JSX, Show, onMount } from 'solid-js'
import { Portal } from 'solid-js/web'

type ModalProps = {
  isOpen: boolean
  onClose: () => void
  title?: string
  children: JSX.Element
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
}

export function Modal(props: ModalProps) {
  const [isAnimating, setIsAnimating] = createSignal(false)

  const sizeClasses = {
    sm: 'max-w-md',
    md: 'max-w-2xl',
    lg: 'max-w-4xl',
    xl: 'max-w-6xl',
    full: 'max-w-7xl',
  }

  const handleClose = () => {
    setIsAnimating(true)
    setTimeout(() => {
      props.onClose()
      setIsAnimating(false)
    }, 200)
  }

  const handleBackdropClick = (e: MouseEvent) => {
    if (e.target === e.currentTarget) {
      handleClose()
    }
  }

  const handleEscape = (e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      handleClose()
    }
  }

  onMount(() => {
    if (props.isOpen) {
      document.addEventListener('keydown', handleEscape)
    }
    return () => {
      document.removeEventListener('keydown', handleEscape)
    }
  })

  return (
    <Show when={props.isOpen}>
      <Portal>
        <div
          class={`fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 transition-opacity duration-200 ${
            isAnimating() ? 'opacity-0' : 'opacity-100'
          }`}
          onClick={handleBackdropClick}
        >
          <div
            class={`bg-white rounded-lg shadow-xl w-full ${
              sizeClasses[props.size || 'lg']
            } max-h-[90vh] flex flex-col transform transition-all duration-200 ${
              isAnimating() ? 'scale-95 opacity-0' : 'scale-100 opacity-100'
            }`}
          >
            {/* Header */}
            <div class='flex items-center justify-between px-6 py-4 border-b border-gray-200'>
              <h2 class='text-xl font-bold text-gray-900'>{props.title || 'Details'}</h2>
              <button
                onClick={handleClose}
                class='text-gray-400 hover:text-gray-600 transition-colors'
                aria-label='Close modal'
              >
                <svg class='w-6 h-6' fill='none' stroke='currentColor' viewBox='0 0 24 24'>
                  <path stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M6 18L18 6M6 6l12 12' />
                </svg>
              </button>
            </div>

            {/* Content */}
            <div class='overflow-y-auto flex-1 px-6 py-4'>{props.children}</div>
          </div>
        </div>
      </Portal>
    </Show>
  )
}
