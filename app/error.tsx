'use client'

export default function ErrorState({ reset }: { reset: () => void }) {
  return <main className="page-shell"><div className="error-state" role="alert"><h1>Inventory unavailable</h1><p>The local API could not load the catalog.</p><button onClick={reset}>Try again</button></div></main>
}
