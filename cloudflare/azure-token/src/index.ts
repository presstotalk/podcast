export interface Env {
  AZURE_CLIENT_ID: string
  AZURE_CLIENT_SECRET: string
  AZURE_REFRESH_TOKEN: string
  AZURE_REDIRECT_URL: string
  AZURE: KVNamespace
}

export default {
	async scheduled(
		controller: ScheduledController,
		env: Env,
		ctx: ExecutionContext
	): Promise<void> {
    await refreshTokens(env)
	},

	async fetch(
		request: Request,
		env: Env,
		ctx: ExecutionContext
	): Promise<Response> {
    await refreshTokens(env)
    return new Response("success")
  },
}

interface Tokens {
  refresh_token: string
  access_token: string
}

async function refreshTokens(env: Env): Promise<void> {
  console.log('refresh tokens')
  const tokens = await getTokens(env)
  await saveTokens(env, tokens)
}

async function getTokens(env: Env): Promise<Tokens> {
  const resp = await fetch('https://login.microsoftonline.com/common/oauth2/v2.0/token', {
    method: 'POST',
    body: `client_id=${env.AZURE_CLIENT_ID}&redirect_uri=${env.AZURE_REDIRECT_URL}&client_secret=${env.AZURE_CLIENT_SECRET}&refresh_token=${env.AZURE_REFRESH_TOKEN}&grant_type=refresh_token`,
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  })

  if (!resp.ok) {
    const msg = await resp.text()
    throw new Error(`failed to refresh tokens: ${msg}`)
  }

  const data = await resp.json<Tokens | null>()
  const { refresh_token, access_token } = data ?? {}
  if (!refresh_token || !access_token) {
    throw new Error(`can't find tokens in the response`)
  }

  return { refresh_token, access_token }
}

async function saveTokens(env: Env, tokens: Tokens): Promise<void> {
  await env.AZURE.put('tokens', JSON.stringify(tokens))
}
