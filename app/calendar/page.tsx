"use client"

import { useState, useEffect } from "react"
import { ClassSidebar, type Course } from "@/components/class-sidebar"
import { WeeklySchedule } from "@/components/weekly-schedule"
import { SelectedCoursesList } from "@/components/selected-courses-list"
import { SidebarProvider } from "@/components/ui/sidebar"

export default function CoursePlanner() {
    // Store quarters and the currently selected quarter
    const [quarters, setQuarters] = useState<number[]>([])
    const [selectedQuarter, setSelectedQuarter] = useState<number | null>(null)

    // Store selected courses by quarter
    const [selectedCoursesByQuarter, setSelectedCoursesByQuarter] = useState<Record<number, Course[]>>({})

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

    return (
        <div className="flex h-screen overflow-hidden">
            <SidebarProvider>
                <div className="w-80 flex-shrink-0">
                    <ClassSidebar
                        selectedCourses={selectedCourses}
                        onCourseToggle={handleCourseToggle}
                        quarters={quarters}
                        selectedQuarter={selectedQuarter}
                        onQuarterChange={handleQuarterChange}
                    />
                </div>

                <div className="flex flex-1 overflow-hidden">
                    <div className="w-64 border-r overflow-y-auto p-4">
                        <SelectedCoursesList
                            courses={selectedCourses}
                            onRemoveCourse={handleCourseToggle}
                            selectedQuarter={selectedQuarter}
                        />
                    </div>

                    <div className="flex-1 overflow-y-auto p-4">
                        <WeeklySchedule
                            courses={selectedCourses}
                            selectedQuarter={selectedQuarter}
                        />
                    </div>
                </div>
            </SidebarProvider>
        </div>
    )
}

