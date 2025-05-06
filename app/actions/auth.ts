"use server"

import { cookies } from "next/headers"

// In a real app, you would use a database and proper password hashing
const MOCK_USERS = [
  { email: "student@example.com", password: "password123" },
  { email: "test@example.com", password: "test123" },
]

export async function loginUser(formData: FormData) {
  const email = formData.get("email") as string
  const password = formData.get("password") as string
  const remember = formData.get("remember") === "on"

  // Validate inputs
  if (!email || !password) {
    return { success: false, error: "Email and password are required" }
  }

  // Check credentials (mock authentication)
  const user = MOCK_USERS.find((user) => user.email === email && user.password === password)

  if (!user) {
    return { success: false, error: "Invalid email or password" }
  }

  // Get cookie store - properly awaited
  const cookieStore = await cookies()
  cookieStore.set("auth-token", email, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    maxAge: remember ? 60 * 60 * 24 * 7 : 60 * 60 * 24, // 7 days or 1 day
    path: "/",
  })

  return { success: true, email }
}

export async function checkAuth() {
  const cookieStore = await cookies()
  const token = cookieStore.get("auth-token")
  return token ? { isLoggedIn: true, email: token.value } : { isLoggedIn: false }
}

export async function logoutUser() {
  const cookieStore = await cookies()
  cookieStore.delete("auth-token")
  return { success: true }
}
