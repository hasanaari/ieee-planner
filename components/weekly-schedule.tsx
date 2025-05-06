"use client"

import { useMemo } from "react"
import type { Course } from "./class-sidebar"
import { getQuarterName } from "@/lib/utils"

type WeeklyScheduleProps = {
    courses: Course[]
    selectedQuarter: number | null
}

type ScheduleSlot = {
    course: Course
    rowStart: number
    rowSpan: number
}

type DaySchedule = {
    [timeSlot: string]: ScheduleSlot[]
}

const DAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"]
const DAY_ABBREVIATIONS = ["M", "T", "W", "Th", "F"]
const TIME_SLOTS = Array.from({ length: 14 }, (_, i) => i + 8) // 8am to 9pm

export function WeeklySchedule({ courses, selectedQuarter }: WeeklyScheduleProps) {
    // Process courses to determine their schedule slots
    const scheduleData = useMemo(() => {
        const daySchedules: Record<string, DaySchedule> = {
            Monday: {},
            Tuesday: {},
            Wednesday: {},
            Thursday: {},
            Friday: {},
        }

        courses.forEach((course) => {
            if (!course.meetingTimes || course.meetingTimes.length === 0) return

            course.meetingTimes.forEach((meeting) => {
                if (meeting.timerange === "TBA" || !meeting.days) return

                // Parse time range (e.g., "10:00 AM - 11:50 AM")
                const timeMatch = meeting.timerange.match(/(\d+):(\d+)\s*(AM|PM)\s*-\s*(\d+):(\d+)\s*(AM|PM)/)
                if (!timeMatch) return

                const [_, startHour, startMin, startAmPm, endHour, endMin, endAmPm] = timeMatch

                // Convert to 24-hour format
                const start24Hour =
                    Number.parseInt(startHour) + (startAmPm === "PM" && Number.parseInt(startHour) !== 12 ? 12 : 0)
                const end24Hour = Number.parseInt(endHour) + (endAmPm === "PM" && Number.parseInt(endHour) !== 12 ? 12 : 0)

                // Calculate row position and span
                const rowStart = (start24Hour - 8) * 4 + Math.floor(Number.parseInt(startMin) / 15) + 2 // +2 for header row
                const endRow = (end24Hour - 8) * 4 + Math.floor(Number.parseInt(endMin) / 15) + 2
                const rowSpan = endRow - rowStart

                // Add to schedule for each day
                meeting.days.forEach((day) => {
                    let fullDay = day
                    if (day === "M") fullDay = "Monday"
                    else if (day === "T") fullDay = "Tuesday"
                    else if (day === "W") fullDay = "Wednesday"
                    else if (day === "Th") fullDay = "Thursday"
                    else if (day === "F") fullDay = "Friday"

                    if (!DAYS.includes(fullDay)) return

                    if (!daySchedules[fullDay][rowStart]) {
                        daySchedules[fullDay][rowStart] = []
                    }

                    daySchedules[fullDay][rowStart].push({
                        course,
                        rowStart,
                        rowSpan,
                    })
                })
            })
        })

        return daySchedules
    }, [courses])

    // Get a color for a course (deterministic based on subject)
    const getCourseColor = (course: Course) => {
        const colors = [
            "bg-blue-100 border-blue-300",
            "bg-green-100 border-green-300",
            "bg-purple-100 border-purple-300",
            "bg-yellow-100 border-yellow-300",
            "bg-red-100 border-red-300",
            "bg-pink-100 border-pink-300",
            "bg-indigo-100 border-indigo-300",
            "bg-teal-100 border-teal-300",
        ]
        const index = course.subject.split("").reduce((acc, char) => acc + char.charCodeAt(0), 0) % colors.length
        return colors[index]
    }

    return (
        <div className="space-y-4">
            <h2 className="text-lg font-semibold">
                Weekly Schedule
                {selectedQuarter && <span className="text-sm font-normal ml-2 text-muted-foreground">Quarter {getQuarterName(selectedQuarter)}</span>}
            </h2>

            {courses.length === 0 ? (
                <p className="text-muted-foreground">No courses selected for this quarter</p>
            ) : (
                <div className="border rounded-md overflow-hidden">
                    <div className="grid grid-cols-6 border-b">
                        <div className="p-2 font-medium text-center border-r">Time</div>
                        {DAYS.map((day) => (
                            <div key={day} className="p-2 font-medium text-center border-r last:border-r-0">
                                {day}
                            </div>
                        ))}
                    </div>

                    <div className="relative">
                        {/* Time slots */}
                        {TIME_SLOTS.map((hour) => (
                            <div key={hour} className="grid grid-cols-6 border-b last:border-b-0">
                                <div className="p-2 text-xs text-center border-r">
                                    {hour % 12 || 12}
                                    {hour >= 12 ? "PM" : "AM"}
                                </div>
                                {DAYS.map((day) => (
                                    <div key={`${day}-${hour}`} className="h-12 border-r last:border-r-0"></div>
                                ))}
                            </div>
                        ))}

                        {/* Course blocks */}
                        {DAYS.map((day, dayIndex) =>
                            Object.entries(scheduleData[day]).map(([rowStart, slots]) =>
                                slots.map((slot, slotIndex) => (
                                    <div
                                        key={`${day}-${slot.course.section}-${slotIndex}`}
                                        className={`absolute border rounded-md p-1 overflow-hidden ${getCourseColor(slot.course)}`}
                                        style={{
                                            gridColumn: `${dayIndex + 2}`,
                                            gridRow: `${slot.rowStart} / span ${slot.rowSpan}`,
                                            top: `${(Number.parseInt(rowStart) - 2) * 12}px`,
                                            left: `${(dayIndex + 1) * (100 / 6)}%`,
                                            width: `${100 / 6 - 1}%`,
                                            height: `${slot.rowSpan * 12}px`,
                                        }}
                                    >
                                        <div className="text-xs font-medium line-clamp-1">{slot.course.title}</div>
                                        <div className="text-xs line-clamp-1">{slot.course.number}</div>
                                    </div>
                                ))
                            )
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}

