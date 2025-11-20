import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { useMutation, useQueryClient } from '@tanstack/solid-query'
import { Layout } from '../components/Layout'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { Select } from '../components/Select'
import { createAircraft } from '../api/aircrafts'
import type { CreateAircraftInput, AircraftStatus } from '../api/types'

export function NewAircraft() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [formData, setFormData] = createSignal<CreateAircraftInput>({
    name: '',
    model: '',
    serial_number: '',
    battery_health: 100,
    status: 1, // Active by default
  })

  const [errors, setErrors] = createSignal<Partial<Record<keyof CreateAircraftInput, string>>>({})

  const createMutation = useMutation(() => ({
    mutationFn: createAircraft,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['aircrafts'] })
      navigate('/aircrafts')
    },
    onError: (error: Error) => {
      alert('Failed to create aircraft: ' + error.message)
    },
  }))

  const handleSubmit = (e: Event) => {
    e.preventDefault()
    const data = formData()
    const newErrors: Partial<Record<keyof CreateAircraftInput, string>> = {}

    if (!data.name) newErrors.name = 'Name is required'
    if (!data.model) newErrors.model = 'Model is required'
    if (!data.serial_number) newErrors.serial_number = 'Serial number is required'
    if (data.battery_health < 0 || data.battery_health > 100) {
      newErrors.battery_health = 'Battery health must be between 0-100'
    }

    setErrors(newErrors)

    if (Object.keys(newErrors).length === 0) {
      createMutation.mutate(data)
    }
  }

  const updateField = <K extends keyof CreateAircraftInput>(field: K, value: CreateAircraftInput[K]) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    setErrors((prev) => ({ ...prev, [field]: undefined }))
  }

  return (
    <Layout>
      <div class='max-w-2xl'>
        <div class='mb-6'>
          <h1 class='text-2xl font-bold text-gray-900'>Add New Aircraft</h1>
          <p class='text-sm text-gray-600 mt-1'>Fill in the details to add a new aircraft to the fleet</p>
        </div>

        <form onSubmit={handleSubmit} class='bg-white rounded-lg border border-gray-200 p-6 space-y-4'>
          <Input
            label='Aircraft Name'
            type='text'
            value={formData().name}
            onInput={(e) => updateField('name', e.currentTarget.value)}
            error={errors().name}
            placeholder='e.g., Sky Hawk 1'
            required
          />

          <Input
            label='Model'
            type='text'
            value={formData().model}
            onInput={(e) => updateField('model', e.currentTarget.value)}
            error={errors().model}
            placeholder='e.g., Mavic 3 Pro'
            required
          />

          <Input
            label='Serial Number'
            type='text'
            value={formData().serial_number}
            onInput={(e) => updateField('serial_number', e.currentTarget.value)}
            error={errors().serial_number}
            placeholder='e.g., SN-12345ABC'
            required
          />

          <Input
            label='Battery Health (%)'
            type='number'
            min='0'
            max='100'
            value={formData().battery_health}
            onInput={(e) => updateField('battery_health', parseInt(e.currentTarget.value) || 0)}
            error={errors().battery_health}
            placeholder='0-100'
            required
          />

          <Select
            label='Status'
            value={formData().status.toString()}
            onChange={(e) => updateField('status', parseInt(e.currentTarget.value) as AircraftStatus)}
            options={[
              { value: '0', label: 'Inactive' },
              { value: '1', label: 'Active' },
              { value: '2', label: 'Maintenance' },
            ]}
          />

          <Input
            label='Last Maintenance Date'
            type='datetime-local'
            value={formData().last_maintenance_date || ''}
            onInput={(e) => updateField('last_maintenance_date', e.currentTarget.value)}
          />

          <Input
            label='Maintenance Due Date'
            type='datetime-local'
            value={formData().maintenance_due_date || ''}
            onInput={(e) => updateField('maintenance_due_date', e.currentTarget.value)}
          />

          <div class='flex gap-3 pt-4'>
            <Button type='submit' disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Creating...' : 'Create Aircraft'}
            </Button>
            <Button type='button' variant='ghost' onClick={() => navigate('/aircrafts')}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </Layout>
  )
}
