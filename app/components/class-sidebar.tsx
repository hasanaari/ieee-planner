"use client"

import { useState } from "react"
import { Check, ChevronDown, Search } from "lucide-react"
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

// Sample class data
const CLASS_DEPARTMENTS = [
  {
    name: "Computer Science",
    color: "bg-blue-500",
    classes: [
      {
        id: "cs101",
        name: "Introduction to Programming",
        professor: "Dr. Smith",
        time: "Mon/Wed 10:00-11:30",
        location: "Tech Hall 101",
      },
      {
        id: "cs201",
        name: "Data Structures",
        professor: "Dr. Johnson",
        time: "Tue/Thu 13:00-14:30",
        location: "Tech Hall 203",
      },
      {
        id: "cs350",
        name: "Algorithms",
        professor: "Dr. Williams",
        time: "Mon/Wed 14:00-15:30",
        location: "Science Center 105",
      },
    ],
  },
  {
    name: "Mathematics",
    color: "bg-green-500",
    classes: [
      {
        id: "math101",
        name: "Calculus I",
        professor: "Dr. Brown",
        time: "Mon/Wed/Fri 9:00-10:00",
        location: "Math Building 302",
      },
      {
        id: "math201",
        name: "Linear Algebra",
        professor: "Dr. Davis",
        time: "Tue/Thu 11:00-12:30",
        location: "Math Building 201",
      },
      {
        id: "math350",
        name: "Differential Equations",
        professor: "Dr. Miller",
        time: "Mon/Wed 15:00-16:30",
        location: "Science Center 210",
      },
    ],
  },
  {
    name: "Physics",
    color: "bg-purple-500",
    classes: [
      {
        id: "phys101",
        name: "Physics I: Mechanics",
        professor: "Dr. Wilson",
        time: "Mon/Wed/Fri 11:00-12:00",
        location: "Physics Hall 101",
      },
      {
        id: "phys201",
        name: "Electricity & Magnetism",
        professor: "Dr. Taylor",
        time: "Tue/Thu 9:00-10:30",
        location: "Physics Hall 205",
      },
    ],
  },
  {
    name: "English",
    color: "bg-yellow-500",
    classes: [
      {
        id: "eng101",
        name: "Composition",
        professor: "Dr. Anderson",
        time: "Mon/Wed 13:00-14:30",
        location: "Humanities 110",
      },
      {
        id: "eng220",
        name: "Creative Writing",
        professor: "Dr. Thomas",
        time: "Fri 13:00-16:00",
        location: "Humanities 210",
      },
    ],
  },
]

export type ClassItem = {
  id: string
  name: string
  professor: string
  time: string
  location: string
  department?: string
  color?: string
}

type ClassSidebarProps = {
  selectedClasses: string[]
  onClassToggle: (classId: string, department: string, color: string) => void
}

export function ClassSidebar({ selectedClasses, onClassToggle }: ClassSidebarProps) {
  const [searchQuery, setSearchQuery] = useState("")

  // Filter classes based on search query
  const filteredDepartments = searchQuery
    ? CLASS_DEPARTMENTS.map((dept) => ({
        ...dept,
        classes: dept.classes.filter(
          (cls) =>
            cls.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            cls.professor.toLowerCase().includes(searchQuery.toLowerCase()),
        ),
      })).filter((dept) => dept.classes.length > 0)
    : CLASS_DEPARTMENTS

  return (
    <Sidebar className="border-r">
      <SidebarHeader className="border-b px-4 py-3">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">My Classes</h2>
          <Badge variant="outline">{selectedClasses.length} Selected</Badge>
        </div>
        <div className="mt-2 relative">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search classes..."
            className="pl-8"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
      </SidebarHeader>

      <SidebarContent>
        {filteredDepartments.map((department) => (
          <Collapsible key={department.name} defaultOpen className="group/collapsible">
            <SidebarGroup>
              <SidebarGroupLabel asChild>
                <CollapsibleTrigger className="flex w-full items-center justify-between">
                  <div className="flex items-center">
                    <div className={`mr-2 h-3 w-3 rounded-full ${department.color}`}></div>
                    {department.name}
                  </div>
                  <ChevronDown className="h-4 w-4 transition-transform group-data-[state=open]/collapsible:rotate-180" />
                </CollapsibleTrigger>
              </SidebarGroupLabel>

              <CollapsibleContent>
                <SidebarGroupContent>
                  <SidebarMenu>
                    {department.classes.map((classItem) => (
                      <SidebarMenuItem key={classItem.id}>
                        <SidebarMenuButton
                          className="flex items-center justify-between"
                          onClick={() => onClassToggle(classItem.id, department.name, department.color)}
                        >
                          <div className="flex flex-col items-start">
                            <span>{classItem.name}</span>
                            <span className="text-xs text-muted-foreground">{classItem.time}</span>
                          </div>
                          <div className="flex h-4 w-4 items-center justify-center">
                            {selectedClasses.includes(classItem.id) && <Check className="h-4 w-4" />}
                          </div>
                        </SidebarMenuButton>
                      </SidebarMenuItem>
                    ))}
                  </SidebarMenu>
                </SidebarGroupContent>
              </CollapsibleContent>
            </SidebarGroup>
          </Collapsible>
        ))}
      </SidebarContent>
    </Sidebar>
  )
}

