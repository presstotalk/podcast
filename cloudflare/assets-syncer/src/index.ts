export interface Env {
  BASE_FOLDER: string
  ASSETS_BUCKET: R2Bucket
  AZURE: KVNamespace
  API_URL: string
}

interface GetFileResult {
  id: string
  folder: string
  file: string
  error?: { code: string }
  ['@microsoft.graph.downloadUrl']: string
}

const DRIVE_API_ENDPOINT = 'https://graph.microsoft.com/v1.0/me/drive'

export default {
	async fetch(
		request: Request,
		env: Env,
		ctx: ExecutionContext
	): Promise<Response> {
    if (request.method !== 'POST') {
      return new Response('method not allowed', { status: 405 })
    }

    const { pathname, destPathanme } = extractPathnameFromRequest(request)
    console.info('sync file: ' + pathname)

    if (pathname === '/feeds/podcast.rss') {
      return handleFeed(env, pathname)
    }

    return handleNormalFile(env, pathname, destPathanme)
	},
}

async function handleNormalFile(env: Env, pathname: string, destPathanme: string): Promise<Response> {
    const accessToken = await getAccessToken(env)
    if (!accessToken) {
      return new Response(`failed to retrieve tokens from database`, { status: 500 })
    }

    const url = genOneDriveUrl(env, pathname)
    const res = await fetchFileData(url, accessToken)
    const fileRes = await fetchOneDriveFile(res)
    return uploadToR2(env, destPathanme, fileRes)
}

function extractPathnameFromRequest(request: Request): { pathname: string; destPathanme: string } {
  const url = new URL(request.url)
  return {
    pathname: url.pathname,
    destPathanme: url.searchParams.get('dest') ?? url.pathname,
  }
}

async function getAccessToken(env: Env): Promise<string | null> {
  const rawTokens = await env.AZURE.get('tokens')
  if (!rawTokens) {
    return null
  }
  const tokens = JSON.parse(rawTokens)
  return tokens?.access_token
}

function genOneDriveUrl(env: Env, pathname: string): string {
  return `${DRIVE_API_ENDPOINT}/root${wrapPathName(env, pathname)}?select=id,etag,folder,file,%40microsoft.graph.downloadUrl&expand=children(select%3Did,etag,folder,file)`;
}

function wrapPathName(env: Env, pathname: string): string {
  const isRequestFolder = pathname.endsWith('/')
  pathname = env.BASE_FOLDER + pathname
  const isIndexingRoot = pathname === '/'
  if (isRequestFolder) {
    if (isIndexingRoot) return ''
    return `:${pathname.replace(/\/$/, '')}:`
  }
  return `:${pathname}`
}

async function fetchFileData(url: string, accessToken: string): Promise<Response> {
  return await fetch(url, {
    headers: {
      Authorization: `bearer ${accessToken}`
    }
  })
}

async function fetchOneDriveFile(res: Response): Promise<Response> {
  const body = await res.json<GetFileResult>()

  if (!res.ok) {
    let status = 500
    if (body.error?.code === 'ItemNotFound') {
      status = 404
    }
    return new Response(JSON.stringify(body.error), { status })
  }

  if (!body.file) {
    return new Response('the resource is not a file', { status: 400 })
  }

  const downloadUrl = body['@microsoft.graph.downloadUrl']
  return fetch(downloadUrl, { method: 'GET' })
}

async function uploadToR2(env: Env, destPathanme: string, res: Response): Promise<Response> {
  if (!res.ok) {
    return res
  }
  await env.ASSETS_BUCKET.put(destPathanme, res.body)
  return new Response('success', { status: 200 })
}

async function handleFeed(env: Env, pathname: string): Promise<Response> {
  const res = await fetch(`${env.API_URL}/feeds/podcast`)
  return uploadToR2(env, pathname, res)
}
