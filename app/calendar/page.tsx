"use client"

import { useState, useEffect, useMemo } from "react"
import { ClassSidebar, type Course } from "@/components/class-sidebar"
import { WeeklySchedule } from "@/components/weekly-schedule"
import { SelectedCoursesList } from "@/components/selected-courses-list"
import { SidebarProvider } from "@/components/ui/sidebar"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ChatSection } from "@/components/chat-section"
import { MajorRequirements } from "@/components/major-requirements"
import { type Message } from '../actions/chat'

export default function CoursePlanner() {
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            content: "Hello! How can I help you with your course planning today?",
            sender: "assistant",
            timestamp: new Date(),
            isError: false
        }
    ])

    const [quarters, setQuarters] = useState<number[]>([])
    const [selectedQuarter, setSelectedQuarter] = useState<number | null>(null)
    const [selectedCoursesByQuarter, setSelectedCoursesByQuarter] = useState<Record<number, Course[]>>({})

    const [majors, setMajors] = useState<string[]>([])
    const [selectedMajor, setSelectedMajor] = useState<string | null>(null)

    // Fetch quarters on component mount
    useEffect(() => {
        const fetchQuarters = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/quarters")
                if (!response.ok) throw new Error("Failed to fetch quarters")
                const data = await response.json()
                setQuarters(data)

                // Auto-select the first quarter
                if (data.length > 0 && !selectedQuarter) {
                    setSelectedQuarter(data[0])
                }
            } catch (error) {
                console.error("Error fetching quarters:", error)
            }
        }

        fetchQuarters()
    }, [])

    // Fetch majors on component mount
    useEffect(() => {
        const fetchMajors = async () => {
            try {
                const response = await fetch("http://localhost:8080/api/majors")
                if (!response.ok) throw new Error("Failed to fetch majors")
                const data = await response.json()
                setMajors(data)
            } catch (error) {
                console.error("Error fetching majors:", error)
            }
        }

        fetchMajors()
    }, [])

    // Get selected courses for the current quarter
    const selectedCourses = selectedQuarter ? selectedCoursesByQuarter[selectedQuarter] || [] : []

    const handleCourseToggle = (course: Course) => {
        if (!selectedQuarter) return

        setSelectedCoursesByQuarter((prev) => {
            // Get current selected courses for this quarter or initialize empty array
            const currentQuarterCourses = prev[selectedQuarter] || []

            // Check if course is already selected
            const isSelected = currentQuarterCourses.some(
                (c) => c.section === course.section && c.quarter === course.quarter
            )

            let updatedQuarterCourses
            if (isSelected) {
                // Remove course if already selected
                updatedQuarterCourses = currentQuarterCourses.filter(
                    (c) => !(c.section === course.section && c.quarter === course.quarter)
                )
            } else {
                // Add course if not selected
                updatedQuarterCourses = [...currentQuarterCourses, course]
            }

            // Return updated map of quarter to selected courses
            return {
                ...prev,
                [selectedQuarter]: updatedQuarterCourses
            }
        })
    }

    const handleQuarterChange = (quarter: number) => {
        setSelectedQuarter(quarter)
    }

    const allSelectedCourses = useMemo(() => {
        const allCourses: string[] = []
        for (const [_, value] of Object.entries(selectedCoursesByQuarter)) {
            for (const course of value) {
                const parts = course.number.split("-")
                const courseKey = (course.subject + " " + parts[0] + "-" + parts[1]).toUpperCase();
                allCourses.push(courseKey);
            }
        }
        return allCourses
    }, [selectedCoursesByQuarter]);

    return (
        <div className="flex h-screen overflow-hidden">
            <SidebarProvider>
                <div className="w-80 flex-shrink-0">
                    <ClassSidebar
                        quarters={quarters}
                        selectedCourses={selectedCourses}
                        onCourseToggle={handleCourseToggle}
                        onQuarterChange={handleQuarterChange}
                        selectedQuarter={selectedQuarter}
                    />
                </div>

                <div className="flex flex-1 overflow-hidden">
                    <div className="w-64 border-r overflow-y-auto p-4">
                        <SelectedCoursesList
                            selectedQuarter={selectedQuarter}
                            courses={selectedCourses}
                            onRemoveCourse={handleCourseToggle} />
                    </div>

                    <div className="flex-1 overflow-y-auto">
                        <Tabs defaultValue="schedule" className="w-full">
                            <TabsList className="grid w-full grid-cols-3">
                                <TabsTrigger value="schedule">Weekly Schedule</TabsTrigger>
                                <TabsTrigger value="chat">Chat</TabsTrigger>
                                <TabsTrigger value="requirements">Major Requirements</TabsTrigger>
                            </TabsList>

                            <TabsContent value="schedule" className="p-4">
                                <WeeklySchedule selectedQuarter={selectedQuarter} courses={selectedCourses} />
                            </TabsContent>

                            <TabsContent value="chat" className="p-4">
                                <ChatSection selectedCourses={allSelectedCourses} allquarters={quarters} selectedMajor={selectedMajor || ""} messages={messages} setMessages={setMessages} />
                            </TabsContent>

                            <TabsContent value="requirements" className="p-4">
                                <MajorRequirements allcourses={allSelectedCourses} majors={majors} selectedMajor={selectedMajor} setSelectedMajor={setSelectedMajor} />
                            </TabsContent>
                        </Tabs>
                    </div>
                </div>
            </SidebarProvider>
        </div>
    )
}
