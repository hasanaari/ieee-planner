import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

export function middleware(request: NextRequest) {
  const authToken = request.cookies.get("auth-token")
  const isAuthPage = request.nextUrl.pathname === "/login"

  // If trying to access protected page without auth
  if (!authToken && !isAuthPage && request.nextUrl.pathname !== "/") {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // If trying to access login page with auth
  if (authToken && isAuthPage) {
    return NextResponse.redirect(new URL("/calendar", request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ["/login", "/calendar/:path*"],
}

