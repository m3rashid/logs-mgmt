import { Button } from './Button'
import type { PaginationType } from '../api/types'

type PaginationProps = {
  pagination: PaginationType
  onPageChange: (page: number) => void
}

export function Pagination(props: PaginationProps) {
  return (
    <div class='flex items-center justify-between px-4 py-3 bg-white border-t border-gray-200'>
      <div class='flex items-center gap-2 text-sm text-gray-700'>
        <span>
          Page {props.pagination.page} of {Math.ceil(props.pagination.count / props.pagination.limit)}
        </span>
        <span class='text-gray-400'>•</span>
        <span>{props.pagination.count} total items</span>
      </div>
      <div class='flex items-center gap-2'>
        <Button
          variant='ghost'
          size='sm'
          onClick={() => props.onPageChange(props.pagination.page - 1)}
          disabled={!props.pagination.hasPreviousPage}
        >
          Previous
        </Button>
        <Button
          variant='ghost'
          size='sm'
          onClick={() => props.onPageChange(props.pagination.page + 1)}
          disabled={!props.pagination.hasNextPage}
        >
          Next
        </Button>
      </div>
    </div>
  )
}
