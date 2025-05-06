"use client"

import { useEffect, useState } from "react"
import { Check, ChevronDown, Search } from 'lucide-react'
import { Input } from "@/components/ui/input"
import { Badge } from "@/components/ui/badge"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
} from "@/components/ui/sidebar"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { getQuarterName } from "@/lib/utils"

// Course type based on the API response
export type Course = {
    title: string
    number: string
    topic: string
    instructors?: {
        name: string
        phone: string
        email: string
        officehours: string
        address: string
    }[]
    meetingTimes: {
        location: string
        days: string[] | null
        starttime: string
        endtime: string
        timerange: string
    }[]
    overview: string
    url: string
    section: number
    subject: string
    school: string
    quarter: number
}

type ClassSidebarProps = {
    selectedCourses: Course[]
    onCourseToggle: (course: Course) => void
    quarters: number[]
    selectedQuarter: number | null
    onQuarterChange: (quarter: number) => void
}

export function ClassSidebar({
    selectedCourses,
    onCourseToggle,
    quarters,
    selectedQuarter,
    onQuarterChange
}: ClassSidebarProps) {
    const [courses, setCourses] = useState<Course[]>([])
    const [searchQuery, setSearchQuery] = useState("")
    const [loading, setLoading] = useState(false)

    // Fetch courses when quarter changes
    useEffect(() => {
        if (!selectedQuarter) return

        const fetchCourses = async () => {
            setLoading(true)
            try {
                const response = await fetch(`http://localhost:8080/api/courses?quarter=${selectedQuarter}`)
                if (!response.ok) throw new Error("Failed to fetch courses")
                const data = await response.json()
                setCourses(data)
            } catch (error) {
                console.error("Error fetching courses:", error)
                setCourses([])
            } finally {
                setLoading(false)
            }
        }

        fetchCourses()
    }, [selectedQuarter])

    // Group courses by school and then by subject
    const coursesBySchoolAndSubject = courses.reduce<Record<string, Record<string, Course[]>>>((acc, course) => {
        if (!acc[course.school]) {
            acc[course.school] = {}
        }

        if (!acc[course.school][course.subject]) {
            acc[course.school][course.subject] = []
        }

        acc[course.school][course.subject].push(course)
        return acc
    }, {})

    // Filter courses based on search query
    const filteredSchools = Object.entries(coursesBySchoolAndSubject)
        .map(([school, subjects]) => {
            const filteredSubjects = Object.entries(subjects)
                .map(([subject, subjectCourses]) => ({
                    subject,
                    courses: subjectCourses.filter((course) => course.title.toLowerCase().includes(searchQuery.toLowerCase()) || course.number.toLowerCase().includes(searchQuery.toLowerCase())),
                }))
                .filter((group) => group.courses.length > 0)

            return {
                school,
                subjects: filteredSubjects,
                hasMatchingCourses: filteredSubjects.length > 0,
            }
        })
        .filter((school) => school.hasMatchingCourses)

    // Get a color for a school (deterministic based on school name)
    const getSchoolColor = (school: string) => {
        const colors = [
            "bg-blue-600",
            "bg-green-600",
            "bg-purple-600",
            "bg-yellow-600",
            "bg-red-600",
            "bg-pink-600",
            "bg-indigo-600",
            "bg-teal-600",
        ]
        const index = school.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0) % colors.length
        return colors[index]
    }

    // Get a color for a subject (deterministic based on subject name)
    const getSubjectColor = (subject: string) => {
        const colors = [
            "bg-blue-500",
            "bg-green-500",
            "bg-purple-500",
            "bg-yellow-500",
            "bg-red-500",
            "bg-pink-500",
            "bg-indigo-500",
            "bg-teal-500",
        ]
        const index = subject.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0) % colors.length
        return colors[index]
    }

    // Check if a course is selected
    const isCourseSelected = (course: Course) => {
        return selectedCourses.some((c) => c.section === course.section && c.quarter === course.quarter)
    }

    return (
        <Sidebar className="border-r h-full">
            <SidebarHeader className="border-b px-4 py-3">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold">Course Selection</h2>
                    <Badge variant="outline">{selectedCourses.length} Selected</Badge>
                </div>

                <div className="mt-2 mb-2">
                    <Select
                        value={selectedQuarter?.toString()}
                        onValueChange={(value) => onQuarterChange(Number.parseInt(value))}
                    >
                        <SelectTrigger className="w-full">
                            <SelectValue placeholder="Select Quarter" />
                        </SelectTrigger>
                        <SelectContent>
                            {quarters.map((quarter) => (
                                <SelectItem key={quarter} value={quarter.toString()}>
                                    {getQuarterName(quarter)}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                <div className="relative">
                    <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                        placeholder="Search courses by title..."
                        className="pl-8"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                </div>
            </SidebarHeader>

            <SidebarContent className="overflow-y-auto">
                {loading ? (
                    <div className="flex justify-center p-4">
                        <p>Loading courses...</p>
                    </div>
                ) : filteredSchools.length === 0 ? (
                    <div className="p-4 text-center text-muted-foreground">
                        {selectedQuarter ? "No courses found" : "Select a quarter to view courses"}
                    </div>
                ) : (
                    filteredSchools.map(({ school, subjects }) => (
                        <Collapsible key={school} className="group/school-collapsible">
                            <SidebarGroup>
                                <SidebarGroupLabel asChild>
                                    <CollapsibleTrigger className="flex w-full items-center justify-between">
                                        <div className="flex items-center">
                                            <div className={`mr-2 h-3 w-3 rounded-full ${getSchoolColor(school)}`}></div>
                                            <span className="font-medium">{school}</span>
                                        </div>
                                        <ChevronDown className="h-4 w-4 transition-transform group-data-[state=open]/school-collapsible:rotate-180" />
                                    </CollapsibleTrigger>
                                </SidebarGroupLabel>

                                <CollapsibleContent>
                                    <SidebarGroupContent>
                                        {subjects.map(({ subject, courses }) => (
                                            <Collapsible key={subject} className="group/subject-collapsible mt-1">
                                                <div className="pl-2">
                                                    <CollapsibleTrigger className="flex w-full items-center justify-between rounded-md px-2 py-1.5 text-sm hover:bg-muted">
                                                        <div className="flex items-center">
                                                            <div className={`mr-2 h-2.5 w-2.5 rounded-full ${getSubjectColor(subject)}`}></div>
                                                            <span>{subject}</span>
                                                        </div>
                                                        <ChevronDown className="h-3.5 w-3.5 transition-transform group-data-[state=open]/subject-collapsible:rotate-180" />
                                                    </CollapsibleTrigger>
                                                </div>

                                                <CollapsibleContent>
                                                    <SidebarMenu className="pl-4">
                                                        {courses.map((course) => (
                                                            <SidebarMenuItem key={`${course.section}-${course.quarter}`}>
                                                                <SidebarMenuButton
                                                                    className="flex items-center justify-between"
                                                                    onClick={() => onCourseToggle(course)}
                                                                >
                                                                    <div className="flex flex-col items-start max-w-[85%]">
                                                                        <span className="font-medium text-sm line-clamp-2">{course.title}</span>
                                                                    </div>
                                                                    <div className="flex h-4 w-4 items-center justify-center">
                                                                        {isCourseSelected(course) && <Check className="h-4 w-4" />}
                                                                    </div>
                                                                </SidebarMenuButton>
                                                            </SidebarMenuItem>
                                                        ))}
                                                    </SidebarMenu>
                                                </CollapsibleContent>
                                            </Collapsible>
                                        ))}
                                    </SidebarGroupContent>
                                </CollapsibleContent>
                            </SidebarGroup>
                        </Collapsible>
                    ))
                )}
            </SidebarContent>
        </Sidebar>
    )
}

