"use client"

import { useState, useEffect } from "react"
import { ChevronLeft, ChevronRight } from "lucide-react"
import { Button } from "@/components/ui/button"

const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]
// contexts

export const Calendar = ({flag = false }) => {
  const [currentDate, setCurrentDate] = useState(new Date())

  const firstDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1)
  const lastDayOfMonth = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0)
  const daysInMonth = lastDayOfMonth.getDate()
  const startingDayOfWeek = firstDayOfMonth.getDay()

  const prevMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() - 1, 1))
  }

  const nextMonth = () => {
    setCurrentDate(new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 1))
  }

  const renderDays = () => {
    const days = []
    for (let i = 0; i < startingDayOfWeek; i++) {
      days.push(<div key={`empty-${i}`} className="h-12 md:h-24"></div>)
    }
    for (let i = 1; i <= daysInMonth; i++) {
      const isToday =
        new Date().toDateString() === new Date(currentDate.getFullYear(), currentDate.getMonth(), i).toDateString()
      days.push(
        <div key={i} className={`h-12 md:h-24 border border-border p-1 ${isToday ? "bg-primary/10" : ""}`}>
          <span className={`text-sm ${isToday ? "font-bold" : ""}`}>{i}</span>
        </div>,
      )
    }
    return days
  }

    useEffect(() => {
        console.log(currentDate);
  }, [currentDate]);

  if (flag) {
  return (
    <div className="w-full max-w-4xl mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">
          {currentDate.toLocaleString("default", { month: "long", year: "numeric" })}
        </h2>
        <div className="flex gap-2">
          <Button variant="outline" size="icon" onClick={prevMonth}>
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <Button variant="outline" size="icon" onClick={nextMonth}>
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-7 gap-1">
        {DAYS.map((day) => (
          <div key={day} className="text-center font-medium py-2">
            {day}
          </div>
        ))}
        {renderDays()}
      </div>
    </div>
  )
  } else {
      return  (
    <div>
    Hello
    </div>
              )
  }
}



export default function CalendarPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <Calendar flag={true}/>
    </main>
  );
}
