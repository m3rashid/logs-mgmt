import { createSignal } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { Layout } from '../components/Layout'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { Select } from '../components/Select'
import { useAircraftsQuery } from '../api/aircrafts'
import { useCreateFlightMutation } from '../api/flights'
import type { CreateFlightInput } from '../api/types'

export function NewFlight() {
  const navigate = useNavigate()

  const [formData, setFormData] = createSignal<CreateFlightInput>({
    aircraft_id: 0,
    source_lat: 0,
    source_lon: 0,
    destination_lat: 0,
    destination_lon: 0,
    start_time: '',
    end_time: '',
  })

  const [errors, setErrors] = createSignal<Partial<Record<keyof CreateFlightInput, string>>>({})

  const aircraftsQuery = useAircraftsQuery(() => ({ page: 1, limit: 100 }))

  const createMutation = useCreateFlightMutation({
    onSuccess: () => navigate('/'),
    onError: (error: Error) => alert('Failed to create flight: ' + error.message),
  })

  const handleSubmit = (e: Event) => {
    e.preventDefault()
    const data = formData()
    const newErrors: Partial<Record<keyof CreateFlightInput, string>> = {}

    if (!data.aircraft_id) newErrors.aircraft_id = 'Aircraft is required'
    if (!data.start_time) newErrors.start_time = 'Start time is required'
    if (!data.end_time) newErrors.end_time = 'End time is required'

    if (data.start_time && data.end_time && data.start_time >= data.end_time) {
      newErrors.end_time = 'End time must be after start time'
    }

    setErrors(newErrors)

    if (Object.keys(newErrors).length === 0) {
      createMutation.mutate(data)
    }
  }

  const updateField = <K extends keyof CreateFlightInput>(field: K, value: CreateFlightInput[K]) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
    setErrors((prev) => ({ ...prev, [field]: undefined }))
  }

  return (
    <Layout>
      <div class='max-w-2xl'>
        <div class='mb-6'>
          <h1 class='text-2xl font-bold text-gray-900'>Schedule New Flight</h1>
          <p class='text-sm text-gray-600 mt-1'>Fill in the details to schedule a new flight</p>
        </div>

        <form onSubmit={handleSubmit} class='bg-white rounded-lg border border-gray-200 p-6 space-y-4'>
          <Select
            label='Aircraft'
            value={formData().aircraft_id.toString()}
            onChange={(e) => updateField('aircraft_id', parseInt(e.currentTarget.value))}
            error={errors().aircraft_id}
            placeholder='Select an aircraft'
            options={
              aircraftsQuery.data?.docs.map((aircraft) => ({
                value: aircraft.id.toString(),
                label: `${aircraft.name} - ${aircraft.model} (${aircraft.serial_number})`,
              })) || []
            }
            required
          />

          <div class='space-y-3'>
            <h3 class='text-sm font-medium text-gray-700'>Source Location</h3>
            <div class='grid grid-cols-2 gap-4'>
              <Input
                label='Latitude'
                type='number'
                step='any'
                value={formData().source_lat}
                onInput={(e) => updateField('source_lat', parseFloat(e.currentTarget.value) || 0)}
                error={errors().source_lat}
                placeholder='e.g., 40.7128'
                required
              />

              <Input
                label='Longitude'
                type='number'
                step='any'
                value={formData().source_lon}
                onInput={(e) => updateField('source_lon', parseFloat(e.currentTarget.value) || 0)}
                error={errors().source_lon}
                placeholder='e.g., -74.0060'
                required
              />
            </div>
          </div>

          <div class='space-y-3'>
            <h3 class='text-sm font-medium text-gray-700'>Destination Location</h3>
            <div class='grid grid-cols-2 gap-4'>
              <Input
                label='Latitude'
                type='number'
                step='any'
                value={formData().destination_lat}
                onInput={(e) => updateField('destination_lat', parseFloat(e.currentTarget.value) || 0)}
                error={errors().destination_lat}
                placeholder='e.g., 34.0522'
                required
              />

              <Input
                label='Longitude'
                type='number'
                step='any'
                value={formData().destination_lon}
                onInput={(e) => updateField('destination_lon', parseFloat(e.currentTarget.value) || 0)}
                error={errors().destination_lon}
                placeholder='e.g., -118.2437'
                required
              />
            </div>
          </div>

          <div class='grid grid-cols-2 gap-4'>
            <Input
              label='Start Time'
              type='datetime-local'
              value={formData().start_time}
              onInput={(e) => updateField('start_time', e.currentTarget.value)}
              error={errors().start_time}
              required
            />

            <Input
              label='End Time'
              type='datetime-local'
              value={formData().end_time}
              onInput={(e) => updateField('end_time', e.currentTarget.value)}
              error={errors().end_time}
              required
            />
          </div>

          <div class='flex gap-3 pt-4'>
            <Button type='submit' disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Creating...' : 'Schedule Flight'}
            </Button>
            <Button type='button' variant='ghost' onClick={() => navigate('/')}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    </Layout>
  )
}
