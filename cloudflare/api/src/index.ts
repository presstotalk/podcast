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
    return getFeed(env, request.method === 'HEAD')
  }
  return new Response('method not allowed', { status: 405 })
}

async function getFeed(env: Env, headOnly: boolean): Promise<Response> {
  const options = {
    status: 200,
    headers: {
      'Content-Type': 'application/rss+xml; charset=utf-8',
      'Cache-Control': 'max-age=3600',
    },
  }
  if (headOnly) {
    return new Response(null, options)
  }
  const feedContent = await env.PRESSTOTALK.get('podcast-feed')
  return new Response(feedContent, options)
}

async function generateFeed(env: Env): Promise<Response> {
  const res = await fetch(`${env.API_URL}/feeds/podcast`)
  const body = await res.text()
  await env.PRESSTOTALK.put('podcast-feed', body)
  return new Response('success', { status: 200 })
}
