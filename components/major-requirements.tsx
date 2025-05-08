"use client";

import type React from "react";

import { useEffect, useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import { Progress } from "./progress";
import { CheckCircle, Circle, AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

type Requirement = {
  between: {
    courses: string[];
  }[];
};

type RequirementGroup = {
  name: string;
  requirementType: number;
  requirements: Requirement[];
};

type MajorRequirements = {
  isEngineering: boolean;
  major: string;
  allreqs: (
    | RequirementGroup
    | { numreqs: number; requirementType: number; type?: number; name?: string }
  )[];
};

function isSubset(subset: string[], superset: string[]) {
  return subset.every((item) => superset.includes(item));
}

export function MajorRequirements({
  allcourses,
  majors,
  selectedMajor,
  setSelectedMajor,
}: {
  allcourses: string[];
  majors: string[];
  selectedMajor: string | null;
  setSelectedMajor: React.Dispatch<React.SetStateAction<string | null>>;
}) {
  const [requirements, setRequirements] = useState<MajorRequirements | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [completionStats, setCompletionStats] = useState({
    completed: 0,
    total: 0,
    percentage: 0,
  });

  // Fetch requirements when major changes
  useEffect(() => {
    if (!selectedMajor) return;

    const fetchRequirements = async () => {
      setLoading(true);
      try {
        const response = await fetch(
          `http://localhost:8080/api/reqs?major=${encodeURIComponent(
            selectedMajor
          )}`
        );
        if (!response.ok) throw new Error("Failed to fetch requirements");
        const data = await response.json();
        setRequirements(data);
      } catch (error) {
        console.error("Error fetching requirements:", error);
        setRequirements(null);
      } finally {
        setLoading(false);
      }
    };

    fetchRequirements();
  }, [selectedMajor]);

  // Calculate completion statistics when requirements change
  useEffect(() => {
    if (!requirements) {
      setCompletionStats({ completed: 0, total: 0, percentage: 0 });
      return;
    }

    let completed = 0;
    let total = 0;

    requirements.allreqs.forEach((reqGroup) => {
      if (!("requirements" in reqGroup)) {
        // Handle special requirement types
        total += 1;
        // We don't have a way to check completion for these yet
        return;
      }

      if (!reqGroup.requirements) return;

      reqGroup.requirements.forEach((req) => {
        if (req.between.length === 0) return;

        total += 1;

        if (req.between.length === 1) {
          // Single option requirement
          if (isSubset(req.between[0].courses, allcourses)) {
            completed += 1;
          }
        } else {
          // Multiple option requirement
          if (
            req.between.some((option) => isSubset(option.courses, allcourses))
          ) {
            completed += 1;
          }
        }
      });
    });

    const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
    setCompletionStats({ completed, total, percentage });
  }, [requirements, allcourses]);

  // Format major name for display (capitalize each word)
  const formatMajorName = (major: string) => {
    return major
      .split(" ")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  };

  // Render requirements for a group
  const renderRequirements = (requirements: Requirement[]) => {
    return requirements
      .filter((req) => {
        // Skip requirements with empty between arrays
        if (req.between.length === 0) return false;

        // For single option requirements, check if courses array is not empty
        if (req.between.length === 1) {
          return req.between[0].courses && req.between[0].courses.length > 0;
        }

        // For multiple options, check if at least one option has courses
        return req.between.some(
          (option) => option.courses && option.courses.length > 0
        );
      })
      .map((req, index) => {
        // Check if this requirement is completed
        const isCompleted =
          req.between.length === 1
            ? isSubset(req.between[0].courses, allcourses)
            : req.between.some((option) =>
                isSubset(option.courses, allcourses)
              );

        return (
          <div
            key={index}
            className={cn(
              "mb-4 pl-4 border-l-2 py-2 rounded-sm transition-colors",
              isCompleted
                ? "border-green-500 bg-green-50/50"
                : "border-muted hover:border-muted-foreground/50"
            )}
          >
            <div className="text-sm">
              {req.between.length > 1 ? (
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    {isCompleted ? (
                      <CheckCircle className="h-5 w-5 text-green-600 flex-shrink-0" />
                    ) : (
                      <Circle className="h-5 w-5 text-muted-foreground flex-shrink-0" />
                    )}
                    <p
                      className={cn(
                        "font-medium",
                        isCompleted ? "text-green-700" : ""
                      )}
                    >
                      Choose one of the following:
                      {isCompleted && (
                        <Badge
                          variant="outline"
                          className="ml-2 bg-green-100 text-green-800 border-green-200"
                        >
                          Completed
                        </Badge>
                      )}
                    </p>
                  </div>
                  <ul className="list-none pl-7 space-y-2">
                    {req.between.map((option, optIndex) => {
                      const optionCompleted = isSubset(
                        option.courses,
                        allcourses
                      );
                      return (
                        <li
                          key={optIndex}
                          className={cn(
                            "flex items-center gap-2 py-1 px-2 rounded",
                            optionCompleted ? "bg-green-100/50" : ""
                          )}
                        >
                          {optionCompleted ? (
                            <CheckCircle className="h-4 w-4 text-green-600 flex-shrink-0" />
                          ) : (
                            <Circle className="h-4 w-4 text-muted-foreground/70 flex-shrink-0" />
                          )}
                          <span
                            className={optionCompleted ? "text-green-700" : ""}
                          >
                            {option.courses.join(" and ")}
                          </span>
                        </li>
                      );
                    })}
                  </ul>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  {isCompleted ? (
                    <CheckCircle className="h-5 w-5 text-green-600 flex-shrink-0" />
                  ) : (
                    <Circle className="h-5 w-5 text-muted-foreground flex-shrink-0" />
                  )}
                  <div>
                    <span
                      className={cn(
                        "font-medium",
                        isCompleted ? "text-green-700" : ""
                      )}
                    >
                      Required: {req.between[0].courses.join(" and ")}
                    </span>
                    {isCompleted && (
                      <Badge
                        variant="outline"
                        className="ml-2 bg-green-100 text-green-800 border-green-200"
                      >
                        Completed
                      </Badge>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        );
      });
  };

  // Get completion status for a requirement group
  const getGroupCompletionStatus = (reqGroup: RequirementGroup) => {
    if (!reqGroup.requirements) return { completed: 0, total: 0 };

    let completed = 0;
    let total = 0;

    reqGroup.requirements.forEach((req) => {
      if (req.between.length === 0) return;

      total += 1;

      if (req.between.length === 1) {
        if (isSubset(req.between[0].courses, allcourses)) {
          completed += 1;
        }
      } else {
        if (
          req.between.some((option) => isSubset(option.courses, allcourses))
        ) {
          completed += 1;
        }
      }
    });

    return { completed, total };
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <h2 className="text-2xl font-bold tracking-tight">
          Major Requirements
        </h2>

        <div className="w-full max-w-xs">
          <Select value={selectedMajor || ""} onValueChange={setSelectedMajor}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select a major" />
            </SelectTrigger>
            <SelectContent>
              {majors.map((major) => (
                <SelectItem key={major} value={major}>
                  {formatMajorName(major)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {loading ? (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <div className="flex flex-col items-center gap-2">
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
              <p className="text-muted-foreground">Loading requirements...</p>
            </div>
          </CardContent>
        </Card>
      ) : requirements ? (
        <div className="space-y-6">
          <Card className="overflow-hidden border-t-4 border-t-primary">
            <CardHeader className="bg-muted/30 pb-2">
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-2">
                <CardTitle className="text-xl">
                  {formatMajorName(requirements.major)}
                </CardTitle>
                <div className="flex items-center gap-2">
                  <div className="text-sm font-medium">
                    {completionStats.completed} of {completionStats.total}{" "}
                    requirements completed
                  </div>
                  <Badge
                    variant={
                      completionStats.percentage >= 75 ? "success" : "secondary"
                    }
                  >
                    {completionStats.percentage}%
                  </Badge>
                </div>
              </div>
              <Progress
                value={completionStats.percentage}
                className="h-2 mt-2"
                indicatorClassName={cn(
                  completionStats.percentage >= 75
                    ? "bg-green-500"
                    : completionStats.percentage >= 50
                    ? "bg-amber-500"
                    : "bg-primary"
                )}
              />
            </CardHeader>
            <CardContent className="pt-6">
              <Accordion type="multiple" className="w-full">
                {requirements.allreqs.map((reqGroup, index) => {
                  // Handle special requirement types (numreqs)
                  if ("numreqs" in reqGroup) {
                    const reqType = reqGroup.type || 0;
                    let reqName = reqGroup.name || "Design and Communications";

                    if (reqType === 1) {
                      reqName = "Social Sciences/Humanities Theme";
                    } else if (reqType === 2) {
                      reqName = "Unrestricted Electives";
                    }

                    let reqText = "";
                    if (reqType === 1) {
                      reqText = `Complete ${reqGroup.numreqs} courses in Social Sciences/Humanities Theme.`;
                    } else if (reqType === 2) {
                      reqText = `Complete ${reqGroup.numreqs} courses as Unrestricted Electives.`;
                    } else {
                      reqText = `Complete ${reqGroup.numreqs} Design and Communications courses.`;
                    }

                    return (
                      <AccordionItem
                        key={index}
                        value={`item-${index}`}
                        className="border-b"
                      >
                        <AccordionTrigger className="py-4 hover:no-underline hover:bg-muted/20 px-2 rounded-sm group">
                          <div className="flex items-center gap-2 text-left">
                            <AlertCircle className="h-5 w-5 text-amber-500 flex-shrink-0" />
                            <div>
                              <span className="font-medium">{reqName}</span>
                              <span className="ml-2 text-muted-foreground text-sm">
                                {reqGroup.numreqs} courses required
                              </span>
                            </div>
                          </div>
                        </AccordionTrigger>
                        <AccordionContent className="px-4 pb-4 pt-2">
                          <div className="bg-muted/20 p-3 rounded-md text-muted-foreground">
                            {reqText}
                          </div>
                        </AccordionContent>
                      </AccordionItem>
                    );
                  }

                  if (!reqGroup.requirements) {
                    return null;
                  }

                  // Handle regular requirement groups
                  const { completed, total } =
                    getGroupCompletionStatus(reqGroup);
                  const allCompleted = completed === total && total > 0;
                  const someCompleted = completed > 0 && completed < total;

                  return (
                    <AccordionItem
                      key={index}
                      value={`item-${index}`}
                      className="border-b"
                    >
                      <AccordionTrigger className="py-4 hover:no-underline hover:bg-muted/20 px-2 rounded-sm group">
                        <div className="flex items-center gap-2 text-left">
                          {allCompleted ? (
                            <CheckCircle className="h-5 w-5 text-green-600 flex-shrink-0" />
                          ) : someCompleted ? (
                            <div className="relative h-5 w-5 flex-shrink-0">
                              <Circle className="h-5 w-5 text-amber-500 absolute" />
                              <div className="absolute inset-0 flex items-center justify-center text-xs font-bold text-amber-500">
                                {completed}
                              </div>
                            </div>
                          ) : (
                            <Circle className="h-5 w-5 text-muted-foreground flex-shrink-0" />
                          )}
                          <div>
                            <span
                              className={cn(
                                "font-medium",
                                allCompleted ? "text-green-700" : ""
                              )}
                            >
                              {reqGroup.name}
                            </span>
                            {total > 0 && (
                              <span className="ml-2 text-muted-foreground text-sm">
                                {completed}/{total} completed
                              </span>
                            )}
                          </div>
                        </div>
                        {allCompleted && (
                          <Badge
                            variant="outline"
                            className="mr-2 bg-green-100 text-green-800 border-green-200"
                          >
                            Complete
                          </Badge>
                        )}
                      </AccordionTrigger>
                      <AccordionContent className="px-2 pb-4 pt-2">
                        {renderRequirements(reqGroup.requirements)}
                      </AccordionContent>
                    </AccordionItem>
                  );
                })}
              </Accordion>
            </CardContent>
          </Card>
        </div>
      ) : selectedMajor ? (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <div className="flex flex-col items-center gap-2 text-muted-foreground">
              <AlertCircle className="h-8 w-8" />
              <p>No requirements found for this major.</p>
            </div>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <div className="text-center text-muted-foreground">
              <p>Select a major to view requirements.</p>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
