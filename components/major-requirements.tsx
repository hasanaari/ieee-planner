"use client";

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
  allreqs: (RequirementGroup | { numreqs: number; requirementType: number })[];
};

function isSubset(subset: string[], superset: string[]) {
  return subset.every((item) => superset.includes(item));
}

// TODO: the COMPLETED action is a little basic
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
      .map((req, index) => (
        <div key={index} className="mb-2 pl-4 border-l-2 border-muted">
          <div className="text-sm">
            {req.between.length > 1 ? (
              <div>
                {req.between.some((option) =>
                  isSubset(option.courses, allcourses)
                ) ? (
                  <>
                    <p className="font-medium mb-1 text-green-800">
                      [COMPLETED] Choose one of the following:
                    </p>
                    <ul className="list-disc pl-5 space-y-1">
                      {req.between.map((option, optIndex) =>
                        isSubset(option.courses, allcourses) ? (
                          <li key={optIndex} className="text-green-800">
                            [COMPLETED] {option.courses.join(" and ")}
                          </li>
                        ) : (
                          <li key={optIndex} className="text-green-800">
                            {option.courses.join(" and ")}
                          </li>
                        )
                      )}
                    </ul>
                  </>
                ) : (
                  <>
                    <p className="font-medium mb-1">
                      Choose one of the following:
                    </p>
                    <ul className="list-disc pl-5 space-y-1">
                      {req.between.map((option, optIndex) => (
                        <li key={optIndex}>{option.courses.join(" and ")}</li>
                      ))}
                    </ul>
                  </>
                )}
              </div>
            ) : (
              <div>
                {isSubset(req.between[0].courses, allcourses) ? (
                  <p className="text-green-800">
                    [COMPLETED] Required: {req.between[0].courses.join(" and ")}
                  </p>
                ) : (
                  <p>Required: {req.between[0].courses.join(" and ")}</p>
                )}
              </div>
            )}
          </div>
        </div>
      ));
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold">Major Requirements</h2>

      <div className="max-w-md">
        <Select value={selectedMajor || ""} onValueChange={setSelectedMajor}>
          <SelectTrigger>
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

      {loading ? (
        <div className="text-center py-8">Loading requirements...</div>
      ) : requirements ? (
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>{formatMajorName(requirements.major)}</CardTitle>
            </CardHeader>
            <CardContent>
              <Accordion type="multiple" className="w-full">
                {requirements.allreqs.map((reqGroup, index) => {
                  // Handle special requirement types (numreqs)
                  if ("numreqs" in reqGroup) {
                    const reqType = (reqGroup as any).type;
                    let reqName =
                      (reqGroup as any).name || "Design and Communications";

                    if (reqType == 1) {
                      reqName = "Social Sciences/Humanities Theme";
                    } else if (reqType == 2) {
                      reqName = "Unrestricted Electives";
                    }

                    let reqText = "";
                    if (reqType == 1) {
                      reqText =
                        "Complete " +
                        reqGroup.numreqs +
                        " courses in Social Sciences/Humanities Theme.";
                    } else if (reqType == 2) {
                      reqText =
                        "Complete " +
                        reqGroup.numreqs +
                        " courses as Unrestricted Electives.";
                    } else {
                      reqText =
                        "Complete " +
                        reqGroup.numreqs +
                        " design and Communications courses.";
                    }

                    return (
                      <AccordionItem key={index} value={`item-${index}`}>
                        <AccordionTrigger>
                          {reqName}: {reqGroup.numreqs} courses required
                        </AccordionTrigger>
                        <AccordionContent>
                          <p className="text-muted-foreground">{reqText}</p>
                        </AccordionContent>
                      </AccordionItem>
                    );
                  }

                  if (!reqGroup.requirements) {
                    return;
                  }

                  // Handle regular requirement groups
                  return (
                    <AccordionItem key={index} value={`item-${index}`}>
                      <AccordionTrigger>{reqGroup.name}</AccordionTrigger>
                      <AccordionContent>
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
        <div className="text-center py-8 text-muted-foreground">
          No requirements found for this major.
        </div>
      ) : (
        <div className="text-center py-8 text-muted-foreground">
          Select a major to view requirements.
        </div>
      )}
    </div>
  );
}
