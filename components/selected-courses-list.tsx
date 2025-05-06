"use client"

import { X } from 'lucide-react'
import { Button } from "@/components/ui/button"
import type { Course } from "./class-sidebar"
import { getQuarterName } from "@/lib/utils"

type SelectedCoursesListProps = {
    courses: Course[]
    onRemoveCourse: (course: Course) => void
    selectedQuarter: number | null
}

export function SelectedCoursesList({ courses, onRemoveCourse, selectedQuarter }: SelectedCoursesListProps) {
    return (
        <div className="space-y-4">
            <h2 className="text-lg font-semibold">
                Selected Courses
            </h2>
            {selectedQuarter && <h2 className="text-sm font-normal text-muted-foreground">Quarter {getQuarterName(selectedQuarter)}</h2>}

            {courses.length === 0 ? (
                <p className="text-muted-foreground text-sm">No courses selected for this quarter</p>
            ) : (
                <div className="space-y-2">
                    {courses.map((course) => (
                        <div key={`${course.section}-${course.quarter}`} className="p-3 rounded-md border relative group">
                            <Button
                                variant="ghost"
                                size="icon"
                                className="absolute top-1 right-1 h-6 w-6 opacity-0 group-hover:opacity-100 transition-opacity"
                                onClick={() => onRemoveCourse(course)}
                            >
                                <X className="h-4 w-4" />
                                <span className="sr-only">Remove course</span>
                            </Button>

                            <h3 className="font-medium text-sm pr-6">{course.title}</h3>
                            <p className="text-xs text-muted-foreground mt-1">{course.number}</p>

                            {course.instructors && course.instructors.length > 0 && (
                                <p className="text-xs text-muted-foreground mt-1">{course.instructors[0].name}</p>
                            )}

                            {course.meetingTimes && course.meetingTimes.length > 0 && (
                                <div className="text-xs text-muted-foreground mt-1">
                                    {course.meetingTimes[0].timerange !== "TBA" ? (
                                        <>
                                            <p>{course.meetingTimes[0].days?.join("/") || "TBA"}</p>
                                            <p>{course.meetingTimes[0].timerange}</p>
                                            <p>{course.meetingTimes[0].location}</p>
                                        </>
                                    ) : (
                                        <p>Time: TBA</p>
                                    )}
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    )
}

