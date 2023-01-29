export interface Env {
  PRESSTOTALK: KVNamespace
  API_URL: string
}

export default {
	async fetch(
		request: Request,
		env: Env,
		ctx: ExecutionContext
	): Promise<Response> {
    const url = new URL(request.url)
    if (url.pathname === '/feeds/podcast') {
      return handleFeed(request, env)
    }
    const apiUrl = new URL(env.API_URL)
    url.host = apiUrl.host
    return fetch(url.toString(), {
      headers: {
        ...request.headers
      },
    })
	},
};

async function handleFeed(request: Request, env: Env): Promise<Response> {
  if (request.method === 'POST') {
    console.info('generate feed')
    return generateFeed(env)
  }
  if (request.method === 'GET' || request.method === 'HEAD') {
    return getFeed(request, env)
  }
  return new Response('method not allowed', { status: 405 })
}

async function getFeed(request: Request, env: Env): Promise<Response> {
  const key = new Request(removeSearch(request.url), request)
  const cachedVal = await caches.default.match(key)
  if (cachedVal) {
    console.info('use cached feed')
    return cachedVal
  }

  console.info('get feed from database')
  const feedData = await env.PRESSTOTALK.get('podcast-feed')
  if (!feedData) {
    return new Response('failed to retrive data from database', { status: 500 })
  }

  const feed = JSON.parse(feedData)
  const headerOnly = request.method === 'HEAD'

  const res = new Response(headerOnly ? null : feed.body, {
    status: 200,
    headers: {
      ...feed.headers,
      'Cache-Control': 'public',
    },
  })

  if (request.method === 'GET') {
    await caches.default.put(key, res.clone())
  }

  return res
}

async function generateFeed(env: Env): Promise<Response> {
  const res = await fetch(`${env.API_URL}/feeds/podcast`)
  const body = await res.text()
  await env.PRESSTOTALK.put('podcast-feed', JSON.stringify({
    body,
    headers: {
      'ETag': res.headers.get('ETag'),
      'Content-Type': res.headers.get('Content-Type'),
      'Content-Length': res.headers.get('Content-Length'),
    },
  }))
  return new Response('success', { status: 200 })
}

function removeSearch(urlStr: string): string {
  const url = new URL(urlStr)
  url.search = ''
  return url.toString()
}
