import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/data-sources/$dataSourceId/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/data-sources/$dataSourceId/"!</div>
}
