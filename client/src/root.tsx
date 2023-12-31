import './root.css'

import {
  Body,
  ErrorBoundary,
  FileRoutes,
  Head,
  Html,
  Meta,
  Routes,
  Scripts,
  Title,
} from 'solid-start'

// @refresh reload
import { Suspense } from 'solid-js'

export default function Root() {
  return (
    <Html lang='en'>
      <Head>
        <Title>Tic Tac Toe</Title>
        <Meta charset='utf-8' />
        <Meta name='viewport' content='width=device-width, initial-scale=1' />
      </Head>
      <Body class='bg-gray-950 caret-transparent text-white'>
        <Suspense>
          <ErrorBoundary>
            <Routes>
              <FileRoutes />
            </Routes>
          </ErrorBoundary>
        </Suspense>
        <Scripts />
      </Body>
    </Html>
  )
}
