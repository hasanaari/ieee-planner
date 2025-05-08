import OpenAI from "openai";

const OPENAI_API_KEY = process.env.NEXT_PUBLIC_OPENAI_API_KEY;
const OPENAI_MODEL = "gpt-4o-mini";
import { getQuarterName } from "@/lib/utils";

export type Message = {
  id: string;
  content: string;
  sender: "user" | "assistant";
  timestamp: Date;
  isError: boolean;
};

// Configuration constants
const MAX_ITEMS = 15; // Maximum number of items to include in a response
const MAX_OVERVIEW_LENGTH = 120; // Maximum length of course overview text
const MAX_TOOL_CALLS = 3; // Maximum number of tool calls to allow
const RESPONSE_FORMAT_TEMPLATE = `
# [Major] Course Recommendations
## [Quarter]

### Core Courses
| Course Code | Course Name | Description |
|-------------|-------------|-------------|
| CODE 101    | Example     | Brief description |

### Elective Options
- **Course 1**: Description
- **Course 2**: Description

*Please verify with your academic advisor.*
`;

/**
 * Gets quarter information in a formatted manner
 */
function getQuarters(allquarters: number[]) {
  return allquarters.map((quarterId) => ({
    quarterId,
    quarterName: getQuarterName(quarterId),
  }));
}

/**
 * Fetches major requirements with error handling and data optimization
 */
async function getMajorRequirements(major: string) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/reqs?major=${encodeURIComponent(major)}`
    );
    if (!response.ok) throw new Error("Failed to fetch requirements");
    const data = await response.json();
    return simplifyMajorRequirements(data);
  } catch (error) {
    console.error("Error fetching major requirements:", error);
    return { error: "Failed to fetch major requirements" };
  }
}

/**
 * Optimizes major requirements data to reduce token usage
 */
function simplifyMajorRequirements(data: any) {
  if (!data) return null;

  const result = {
    major: data.Major,
    isEngineering: data.IsEngineering,
    requirements: [],
  };

  if (data.AllRequirements && Array.isArray(data.AllRequirements)) {
    for (const req of data.AllRequirements.slice(0, 8)) {
      if (req.RequirementType === 0) {
        // Generic requirements
        const simplified = {
          type: "Generic",
          name: req.Name,
          courses: [],
        };

        // Process requirements
        if (req.Requirements && Array.isArray(req.Requirements)) {
          const courseSet = new Set();

          // Limit the number of requirements processed
          for (const option of req.Requirements.slice(0, 4)) {
            if (option.Between && Array.isArray(option.Between)) {
              for (const between of option.Between.slice(0, 2)) {
                if (between.Courses && Array.isArray(between.Courses)) {
                  between.Courses.slice(0, 3).forEach((course: string) => {
                    courseSet.add(course);
                  });
                }
              }
            }
          }

          simplified.courses = Array.from(courseSet);
        }

        result.requirements.push(simplified);
      } else if (req.RequirementType === 1) {
        // Theme requirements
        result.requirements.push({
          type: "Theme",
          numReqs: req.NumRequirements,
        });
      } else if (req.RequirementType === 2) {
        // Unrestricted requirements
        result.requirements.push({
          type: "Unrestricted",
          numReqs: req.NumRequirements,
        });
      } else {
        // Unknown requirements
        result.requirements.push({
          type: "Unknown",
          numReqs: req.NumRequirements,
        });
      }
    }
  }

  return result;
}

/**
 * Fetches courses by quarter with pagination and data optimization
 */
async function getCoursesByQuarter(quarter: number, limit = MAX_ITEMS) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/courses?quarter=${quarter}&limit=${limit}`
    );
    if (!response.ok) throw new Error("Failed to fetch courses by quarter");
    const data = await response.json();
    return simplifyCoursesData(data);
  } catch (error) {
    console.error("Error fetching courses by quarter:", error);
    return [];
  }
}

/**
 * Fetches courses by subject with data optimization
 */
async function getCoursesBySubject(subject: string, limit = MAX_ITEMS) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/courses/subject?subject=${subject}&limit=${limit}`
    );
    if (!response.ok) throw new Error("Failed to fetch courses by subject");
    const data = await response.json();
    return simplifyCoursesData(data);
  } catch (error) {
    console.error("Error fetching courses by subject:", error);
    return [];
  }
}

/**
 * Fetches courses by key with data optimization
 */
async function getCoursesByKey(key: string) {
  try {
    const response = await fetch(
      `http://localhost:8080/api/courses/key?key=${key}`
    );
    if (!response.ok) throw new Error("Failed to fetch courses by key");
    const data = await response.json();
    return simplifyCoursesData(data);
  } catch (error) {
    console.error("Error fetching courses by key:", error);
    return [];
  }
}

/**
 * Optimizes course data to reduce token usage
 */
function simplifyCoursesData(courses: any[]) {
  if (!Array.isArray(courses)) return [];

  return courses.slice(0, MAX_ITEMS).map((course) => ({
    title: course.Title || "",
    number: course.Number || "",
    subject: course.Subject || "",
    overview: course.Overview
      ? course.Overview.substring(0, MAX_OVERVIEW_LENGTH) +
        (course.Overview.length > MAX_OVERVIEW_LENGTH ? "..." : "")
      : "",
    instructors: Array.isArray(course.Instructors)
      ? course.Instructors.slice(0, 1)
      : [],
    quarter: course.Quarter,
  }));
}

/**
 * Creates a system prompt for clean responses with course formatting
 */
function getSystemPrompt(selectedCourses, selectedMajor, allquarters) {
  return `You are a helpful course planning assistant at Northwestern University.
      
    Please respond in a natural, conversational way - just like a helpful advisor would talk.
    
    Format your recommendations as follows:
    - Start with a friendly greeting
    - Don't use any markdown formatting (no # symbols) or HTML tags
    - Format course codes in ALL CAPS like this: "CS 101: INTRODUCTION TO PROGRAMMING"
    - Mention 3-5 core courses the student should take based on their major
    - Suggest 2-3 interesting electives that align with typical interests for this major
    - End with a friendly reminder to check with an academic advisor
    
    The student has taken: ${selectedCourses.join(", ")}
    Their major is: ${selectedMajor}`;
}

function formatCourseResponse(response: string): string {
  let formatted = response
    .trim()
    .replace(/```[\s\S]*?```/g, "")
    .trim();

  formatted = formatted.replace(/^# Course Recommendations\s*/i, "");

  formatted = formatted.replace(/<[^>]+>/g, "");

  // Remove markdown formatting
  formatted = formatted.replace(/\*\*([^*]+)\*\*/g, "$1"); // Remove bold markdown
  formatted = formatted.replace(/\*([^*]+)\*/g, "$1"); // Remove italic markdown

  // Remove bullet points and numbering
  formatted = formatted
    .replace(/^\s*[-*]\s*/gm, "") // Remove bullet points
    .replace(/^\s*\d+\.\s+/gm, "") // Remove numbered lists
    .replace(/\n\s*---\s*\n/g, "\n\n") // Remove horizontal rules
    .replace(/^#+\s+/gm, "") // Remove header markers
    .replace(/\n{3,}/g, "\n\n") // No more than 2 consecutive newlines
    .trim();

  return formatted;
}
// /**
//  * Extracts course information from a section
//  */
// function extractCoursesFromSection(lines, startIndex, endIndex) {
//   const courses = [];
//   let currentCourse = null;

//   for (let i = startIndex; i < endIndex; i++) {
//     const line = lines[i].trim();

//     // Skip empty lines and section dividers
//     if (line === "" || line === "---") continue;

//     // Skip if this is an advisor note line
//     if (line.includes("advisor")) continue;

//     // Check if this line contains a course code
//     const courseMatch = line.match(/^\*\*(.+?)\*\*\s*(.+)?$/);

//     if (courseMatch) {
//       // Found a new course
//       currentCourse = {
//         code: courseMatch[1].trim(),
//         desc: (courseMatch[2] || "").trim(),
//       };
//       courses.push(currentCourse);
//     } else if (currentCourse && line.length > 0) {
//       // This is a continuation of the previous course description
//       currentCourse.desc += " " + line;
//     }
//   }

//   return courses;
// }

// Initialize OpenAI client
const openAI = new OpenAI({
  apiKey: OPENAI_API_KEY,
  dangerouslyAllowBrowser: true,
});

/**
 * Main chat function with token optimization
 */
export async function askChat(
  selectedCourses: string[],
  selectedMajor: string,
  allquarters: number[],
  userprompt: string
): Promise<string> {
  // Initial context with enhanced system prompt
  const initialContext: OpenAI.Chat.ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: getSystemPrompt(selectedCourses, selectedMajor, allquarters),
    },
    {
      role: "user",
      content: userprompt,
    },
  ];

  // Define tools for OpenAI
  const tools = [
    {
      type: "function",
      function: {
        name: "getMajorRequirements",
        description:
          "Returns the requirements for a major in simplified format.",
        parameters: {
          type: "object",
          properties: {
            major: {
              type: "string",
              description:
                "The major name (e.g., 'Computer Science', 'Civil Engineering')",
            },
          },
          required: ["major"],
        },
      },
    },
    {
      type: "function",
      function: {
        name: "getCoursesByQuarter",
        description:
          "Returns courses available in a specific quarter (limited to most relevant).",
        parameters: {
          type: "object",
          properties: {
            quarterId: {
              type: "integer",
              description: "The quarter ID",
            },
          },
          required: ["quarterId"],
        },
      },
    },
    {
      type: "function",
      function: {
        name: "getCoursesBySubject",
        description:
          "Returns courses for a given subject (limited to most relevant).",
        parameters: {
          type: "object",
          properties: {
            subject: {
              type: "string",
              description: "The subject code (e.g., 'COMP_SCI', 'ECON')",
            },
          },
          required: ["subject"],
        },
      },
    },
    {
      type: "function",
      function: {
        name: "getCoursesByKey",
        description: "Returns details for a specific course key.",
        parameters: {
          type: "object",
          properties: {
            key: {
              type: "string",
              description: "The course key (e.g., 'COMP_SCI 213-0')",
            },
          },
          required: ["key"],
        },
      },
    },
  ];

  let currentContext = [...initialContext];
  let response: OpenAI.Chat.Completions.ChatCompletion;
  let toolCallCount = 0;

  try {
    // Initial request
    response = await openAI.chat.completions.create({
      model: OPENAI_MODEL,
      messages: currentContext,
      tools,
      tool_choice: "auto",
      temperature: 0.7, // Add some variation to responses
    });

    // Handle tool calls loop with token optimization
    while (
      response.choices[0].finish_reason === "tool_calls" &&
      toolCallCount < MAX_TOOL_CALLS
    ) {
      toolCallCount++;
      const assistantMessage = response.choices[0].message;
      const toolResponses: OpenAI.Chat.ChatCompletionMessageParam[] = [];

      // Process tool calls
      for (const toolCall of assistantMessage.tool_calls || []) {
        let result;
        try {
          const args = JSON.parse(toolCall.function.arguments);

          switch (toolCall.function.name) {
            case "getMajorRequirements":
              result = await getMajorRequirements(args.major);
              break;
            case "getCoursesByQuarter":
              result = await getCoursesByQuarter(args.quarterId);
              break;
            case "getCoursesBySubject":
              result = await getCoursesBySubject(args.subject);
              break;
            case "getCoursesByKey":
              result = await getCoursesByKey(args.key);
              break;
            case "getCoursesByQuarter":
              result = await getCoursesByQuarter(args.quarterId);
              break;
            case "getCoursesBySubject":
              result = await getCoursesBySubject(args.subject);
              break;
            case "getCoursesByKey":
              result = await getCoursesByKey(args.key);
              break;
            default:
              throw new Error(`Unknown tool: ${toolCall.function.name}`);
          }

          toolResponses.push({
            role: "tool",
            content: JSON.stringify(result),
            tool_call_id: toolCall.id,
          });
        } catch (error) {
          console.error(
            `Tool call error for ${toolCall.function.name}:`,
            error
          );
          toolResponses.push({
            role: "tool",
            content: JSON.stringify({
              error: `Error executing ${toolCall.function.name}: ${error.message}`,
            }),
            tool_call_id: toolCall.id,
          });
        }
      }

      // Reset context to only essential messages to save tokens
      currentContext = [
        initialContext[0], // System message
        initialContext[1], // User's original question
        assistantMessage, // Assistant's tool call
        ...toolResponses, // Tool responses
      ];

      // Get next response
      response = await openAI.chat.completions.create({
        model: OPENAI_MODEL,
        messages: currentContext,
      });
    }

    // Apply formatting improvements to the final response
    const content = response.choices[0].message.content || "";
    return formatCourseResponse(content);
  } catch (error) {
    console.error("OpenAI API error:", error);
    if (error.status === 429) {
      return "I encountered a limit while processing your request. Please try a more specific question or try again later.";
    }
    return "Sorry, I encountered an error while processing your request. Please try again.";
  }
}

// Export other functions for testing/debugging
export {
  simplifyCoursesData,
  simplifyMajorRequirements,
  formatCourseResponse,
  getSystemPrompt,
};
