export interface Env {
  BASE_FOLDER: string
  AZURE: KVNamespace
}

interface GetFileResult {
  name: string
  eTag: string
  size: number
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
    const pathname = extractPathnameFromRequest(request)

    const accessToken = await getAccessToken(env)
    if (!accessToken) {
      return new Response(`failed to retrieve tokens from database`, { status: 500 })
    }

    const url = genOneDriveUrl(env, pathname)
    const res = await fetchFileData(url, accessToken)
    return fetchOneDriveFile(res)
	},
}

function extractPathnameFromRequest(request: Request): string {
  return decodeURIComponent(new URL(request.url).pathname).toLowerCase()
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
  return `${DRIVE_API_ENDPOINT}/root${wrapPathName(env, pathname)}?select=name,eTag,size,id,folder,file,%40microsoft.graph.downloadUrl&expand=children(select%3Dname,eTag,size,id,folder,file)`;
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
  const remoteRes = await fetch(downloadUrl)
  const { readable, writable } = new TransformStream()
  remoteRes.body?.pipeTo(writable)
  return new Response(readable, remoteRes)
}
