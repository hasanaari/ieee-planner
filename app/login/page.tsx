"use client"

import { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { BookOpen, Clock } from "lucide-react"
import { loginUser } from "../actions/auth"

export default function LoginPage() {
  const router = useRouter()
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  async function handleSubmit(formData: FormData) {
    setIsLoading(true)
    setError(null)

    try {
      const result = await loginUser(formData)

      if (result.success) {
        router.push("/calendar")
      } else {
        setError(result.error || "Login failed")
      }
    } catch (err) {
      setError("An unexpected error occurred")
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-indigo-50 to-purple-50 px-4 py-12 sm:px-6 lg:px-8">
      <Card className="w-full max-w-md border-indigo-100 shadow-lg">
        <CardHeader className="space-y-1">
          <div className="flex justify-center mb-2">
            <BookOpen className="h-12 w-12 text-indigo-600" />
          </div>
          <CardTitle className="text-2xl font-bold text-center text-indigo-800">Study Planner</CardTitle>
          <CardDescription className="text-center">
            Log in to access your study schedule and track your progress
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form action={handleSubmit}>
            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  placeholder="name@example.com"
                  required
                  className="border-indigo-200 focus:border-indigo-400"
                />
              </div>
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="password">Password</Label>
                  <Link
                    href="/forgot-password"
                    className="text-sm font-medium text-indigo-600 hover:text-indigo-800 hover:underline"
                  >
                    Forgot password?
                  </Link>
                </div>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  required
                  className="border-indigo-200 focus:border-indigo-400"
                />
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox id="remember" name="remember" />
                <Label
                  htmlFor="remember"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Remember me
                </Label>
              </div>
            </div>
            <div className="mt-6 flex flex-col space-y-4">
              <Button className="w-full bg-indigo-600 hover:bg-indigo-700" type="submit" disabled={isLoading}>
                {isLoading ? "Signing in..." : "Sign in"}
              </Button>
              <div className="text-center text-sm">
                Don't have an account?{" "}
                <Link href="/register" className="font-medium text-indigo-600 hover:text-indigo-800 hover:underline">
                  Sign up
                </Link>
              </div>
            </div>
          </form>
        </CardContent>
        <CardFooter className="flex justify-center border-t border-indigo-100 pt-4">
          <div className="flex items-center space-x-2 text-sm text-gray-500">
            <Clock className="h-4 w-4" />
            <span>Plan your future</span>
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}

