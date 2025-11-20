import { type JSX } from 'solid-js'
import { A } from '@solidjs/router'

export function Layout(props: Readonly<{ children: JSX.Element }>) {
  return (
    <div class='min-h-screen bg-gray-50'>
      <nav class='bg-white border-b border-gray-200'>
        <div class='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8'>
          <div class='flex justify-between h-14'>
            <div class='flex items-center gap-6'>
              <A href='/' class='text-lg font-bold text-gray-900'>
                Flight Management
              </A>
              <div class='flex gap-4'>
                <NavLink href='/'>Flights</NavLink>
                <NavLink href='/aircrafts'>Aircrafts</NavLink>
              </div>
            </div>
          </div>
        </div>
      </nav>
      <main class='max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6'>{props.children}</main>
    </div>
  )
}

function NavLink(props: { href: string; children: JSX.Element }) {
  return (
    <A
      href={props.href}
      class='text-sm text-gray-600 hover:text-gray-900 transition-colors'
      activeClass='text-blue-600 font-medium'
    >
      {props.children}
    </A>
  )
}
