"use client"

import { Badge } from "@/components/ui/badge"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Calendar } from "@/components/ui/calendar"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { BookOpen, CalendarIcon, Clock, GraduationCap, LogOut, MapPin, User } from "lucide-react"
import { ClassSidebar, type ClassItem } from "../components/class-sidebar"
import { checkAuth, logoutUser } from "../actions/auth"
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"

// Sample class data with full details
const CLASS_DATA: Record<string, ClassItem> = {
  cs101: {
    id: "cs101",
    name: "Introduction to Programming",
    professor: "Dr. Smith",
    time: "Mon/Wed 10:00-11:30",
    location: "Tech Hall 101",
    department: "Computer Science",
    color: "bg-blue-500",
  },
  cs201: {
    id: "cs201",
    name: "Data Structures",
    professor: "Dr. Johnson",
    time: "Tue/Thu 13:00-14:30",
    location: "Tech Hall 203",
    department: "Computer Science",
    color: "bg-blue-500",
  },
  cs350: {
    id: "cs350",
    name: "Algorithms",
    professor: "Dr. Williams",
    time: "Mon/Wed 14:00-15:30",
    location: "Science Center 105",
    department: "Computer Science",
    color: "bg-blue-500",
  },
  math101: {
    id: "math101",
    name: "Calculus I",
    professor: "Dr. Brown",
    time: "Mon/Wed/Fri 9:00-10:00",
    location: "Math Building 302",
    department: "Mathematics",
    color: "bg-green-500",
  },
  math201: {
    id: "math201",
    name: "Linear Algebra",
    professor: "Dr. Davis",
    time: "Tue/Thu 11:00-12:30",
    location: "Math Building 201",
    department: "Mathematics",
    color: "bg-green-500",
  },
  math350: {
    id: "math350",
    name: "Differential Equations",
    professor: "Dr. Miller",
    time: "Mon/Wed 15:00-16:30",
    location: "Science Center 210",
    department: "Mathematics",
    color: "bg-green-500",
  },
  phys101: {
    id: "phys101",
    name: "Physics I: Mechanics",
    professor: "Dr. Wilson",
    time: "Mon/Wed/Fri 11:00-12:00",
    location: "Physics Hall 101",
    department: "Physics",
    color: "bg-purple-500",
  },
  phys201: {
    id: "phys201",
    name: "Electricity & Magnetism",
    professor: "Dr. Taylor",
    time: "Tue/Thu 9:00-10:30",
    location: "Physics Hall 205",
    department: "Physics",
    color: "bg-purple-500",
  },
  eng101: {
    id: "eng101",
    name: "Composition",
    professor: "Dr. Anderson",
    time: "Mon/Wed 13:00-14:30",
    location: "Humanities 110",
    department: "English",
    color: "bg-yellow-500",
  },
  eng220: {
    id: "eng220",
    name: "Creative Writing",
    professor: "Dr. Thomas",
    time: "Fri 13:00-16:00",
    location: "Humanities 210",
    department: "English",
    color: "bg-yellow-500",
  },
}

export default function CalendarPage() {
  const router = useRouter()
  const [date, setDate] = useState<Date | undefined>(new Date())
  const [user, setUser] = useState<{ email: string } | null>(null)
  const [selectedClasses, setSelectedClasses] = useState<string[]>(["cs101", "math101", "phys101"]) // Default selected classes
  const [view, setView] = useState<"week" | "month">("week")

  useEffect(() => {
    async function validateAuth() {
      const auth = await checkAuth()
      if (!auth.isLoggedIn) {
        router.push("/login")
      } else {
        setUser({ email: auth.email })
      }
    }

    validateAuth()
  }, [router])

  async function handleLogout() {
    await logoutUser()
    router.push("/login")
  }

  function handleClassToggle(classId: string, department: string, color: string) {
    setSelectedClasses((prev) => (prev.includes(classId) ? prev.filter((id) => id !== classId) : [...prev, classId]))
  }

  if (!user) {
    return <div className="flex min-h-screen items-center justify-center">Loading...</div>
  }

  // Get selected class details
  const selectedClassDetails = selectedClasses.map((id) => CLASS_DATA[id]).filter(Boolean)

  return (
    <SidebarProvider>
      <div className="min-h-screen bg-background">
        <header className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
            <div className="flex items-center space-x-2">
              <SidebarTrigger className="md:hidden" />
              <GraduationCap className="h-6 w-6 text-indigo-600" />
              <h1 className="text-xl font-bold text-indigo-800">Class Schedule</h1>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600 hidden sm:inline-block">{user.email}</span>
              <Button variant="outline" size="sm" onClick={handleLogout}>
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>
        </header>

        <div className="flex h-[calc(100vh-65px)]">
          <ClassSidebar selectedClasses={selectedClasses} onClassToggle={handleClassToggle} />

          <main className="flex-1 overflow-auto p-4">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-2xl font-bold">My Class Schedule</h2>
              <Tabs value={view} onValueChange={(v) => setView(v as "week" | "month")}>
                <TabsList>
                  <TabsTrigger value="week">Week</TabsTrigger>
                  <TabsTrigger value="month">Month</TabsTrigger>
                </TabsList>
              </Tabs>
            </div>

            <div className="grid grid-cols-1 gap-6">
              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="flex items-center">
                    <CalendarIcon className="h-5 w-5 mr-2 text-indigo-600" />
                    {view === "week" ? "Weekly Schedule" : "Monthly Calendar"}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Calendar mode="single" selected={date} onSelect={setDate} className="rounded-md border" />
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="pb-2">
                  <CardTitle className="flex items-center">
                    <Clock className="h-5 w-5 mr-2 text-indigo-600" />
                    Classes for {date?.toLocaleDateString("en-US", { weekday: "long", month: "long", day: "numeric" })}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  {selectedClassDetails.length > 0 ? (
                    <div className="space-y-3">
                      {selectedClassDetails.map((classItem) => (
                        <div key={classItem.id} className="flex items-start space-x-3 p-3 rounded-lg border">
                          <div className={`${classItem.color} h-10 w-1 rounded-full flex-shrink-0`}></div>
                          <div className="flex-1">
                            <h3 className="font-medium">{classItem.name}</h3>
                            <div className="mt-1 grid grid-cols-1 gap-1 text-sm text-muted-foreground">
                              <div className="flex items-center">
                                <User className="mr-2 h-3.5 w-3.5" />
                                {classItem.professor}
                              </div>
                              <div className="flex items-center">
                                <Clock className="mr-2 h-3.5 w-3.5" />
                                {classItem.time}
                              </div>
                              <div className="flex items-center">
                                <MapPin className="mr-2 h-3.5 w-3.5" />
                                {classItem.location}
                              </div>
                            </div>
                          </div>
                          <div className="flex-shrink-0">
                            <Badge variant="outline">{classItem.department}</Badge>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
                      <BookOpen className="h-12 w-12 mb-2 opacity-20" />
                      <p>No classes selected for this day</p>
                      <p className="text-sm">Select classes from the sidebar to view your schedule</p>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          </main>
        </div>
      </div>
    </SidebarProvider>
  )
}

