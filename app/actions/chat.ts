import OpenAI from "openai";

const OPENAI_API_KEY = process.env.NEXT_PUBLIC_OPENAI_API_KEY;
const OPENAI_MODEL = 'gpt-4o-mini';
import { getQuarterName } from "@/lib/utils"

export type Message = {
    id: string
    content: string
    sender: "user" | "assistant"
    timestamp: Date
    isError: boolean
}

// Get quarters -> returns all quarters with their names (i.e. the dates) and ids
function getQuarters(allquarters: number[]) {
    return allquarters.map((quarterId) => ({
        quarterId,
        quarterName: getQuarterName(quarterId)
    }));
}
// Get taken courses -> returns keys of all taken courses
// Get selected major -> returns selected major(id is string same as name)

// Get major requirements -> returns a major’s requirements
async function getMajorRequirements(major: string) {
    try {
        const response = await fetch(`http://localhost:8080/api/reqs?major=${encodeURIComponent(major)}`)
        if (!response.ok) throw new Error("Failed to fetch requirements")
        const data = await response.json()
        return data;
    } catch (error) {
        return null;
    }
}
// Get courses by quarter -> returns all courses by quarter id
async function getCoursesByQuarter(quarter: number) {
    try {
        const response = await fetch(`http://localhost:8080/api/courses?quarter=${quarter}`)
        if (!response.ok) throw new Error("Failed to fetch courses by quarter")
        const data = await response.json()
        return data
    } catch (error) {
        return []
    }
}

// Get courses by subject -> returns all courses that belong to a certain subject
async function getCoursesBySubject(subject: string) {
    try {
        const response = await fetch(`http://localhost:8080/api/courses/subject?subject=${subject}`)
        if (!response.ok) throw new Error("Failed to fetch courses by subject")
        const data = await response.json()
        return data
    } catch (error) {
        return []
    }
}
// Get courses by key -> returns all course listings within available quarters
async function getCoursesByKey(key: string) {
    try {
        const response = await fetch(`http://localhost:8080/api/courses/key?key=${key}`)
        if (!response.ok) throw new Error("Failed to fetch courses by key")
        const data = await response.json()
        return data
    } catch (error) {
        return []
    }
}

function getSystemPrompt(selectedCourses: string[], selectedMajor: string, allquarters: number[]) {
    return 'You are an assistant for course planning that gives information about courses and gives suggestions about courses to a studnet at Northwestern University.\n' +
        "The student is interested in the major " + selectedMajor + ".\n" +
        "The student has currently taken or is taking the following courses(keys) " + selectedCourses.join(",") + ".\n" +
        "The available quarters in the system are " + JSON.stringify(getQuarters(allquarters)) + "."
}

const openAI = new OpenAI({
    apiKey: OPENAI_API_KEY,
    dangerouslyAllowBrowser: true
});

async function fulfillFunctionCalls(response: OpenAI.Chat.Completions.ChatCompletion, context: OpenAI.Chat.Completions.ChatCompletionMessageParam[]): Promise<boolean> {
    const willInvokeFunction = response.choices[0].finish_reason == 'tool_calls'
    if (!willInvokeFunction) return false;
    for (const toolcall of response.choices[0].message.tool_calls!) {
        const toolname = toolcall.function.name;
        let callres: string;
        const rawArgument = toolcall.function.arguments;
        const parsedArguments = JSON.parse(rawArgument);

        switch (toolname) {
            case 'getCoursesBySubject':
                callres = JSON.stringify(await getCoursesBySubject(parsedArguments.subject));
                break;
            case 'getCoursesByQuarter':
                callres = JSON.stringify(await getCoursesByQuarter(parsedArguments.quarterId));
                break;
            case 'getMajorRequirements':
                callres = JSON.stringify(await getMajorRequirements(parsedArguments.major));
                break
            case 'getCoursesByKey':
                callres = JSON.stringify(await getCoursesByKey(parsedArguments.key));
                break;
            default:
                callres = ""
        }

        context.push(response.choices[0].message);
        context.push({
            role: 'tool',
            content: callres,
            tool_call_id: toolcall.id
        })
    }
    return true;
}

// TODO: probably more capabilities as tools and one thing to make sure this guy doesnt mess up tool calls is to provide all subjects
export async function askChat(selectedCourses: string[], selectedMajor: string, allquarters: number[], userprompt: string): Promise<string> {
    const context: OpenAI.Chat.ChatCompletionMessageParam[] = [
        {
            role: 'system',
            content: getSystemPrompt(selectedCourses, selectedMajor, allquarters)
        },
        {
            role: 'user',
            content: userprompt
        }
    ]

    let response = await openAI.chat.completions.create({
        model: OPENAI_MODEL,
        messages: context,
        tools: [
            {
                type: 'function',
                function: {
                    name: 'getMajorRequirements',
                    description: 'Returns the requirements for a major. Majors are stored as MajorRequirements objects containing a name, engineering flag, and a list of mixed-type requirement structs (Generic, Theme, Unrestricted, Unknown). The requirement structs contain a list of options, options each contain a list of requirement, requirements each can contain 1 or more course keys.',
                    parameters: {
                        type: 'object',
                        properties: {
                            major: {
                                type: 'string',
                                description: 'The major for which the requirements will be fetched'
                            }
                        },
                        required: ['major']
                    }
                }
            },
            {
                type: 'function',
                function: {
                    name: 'getCoursesByQuarter',
                    description: 'Returns the list of all courses in a quarter. This will be a lot of courses. The courses themselves have a Title, a Number, a Topic, Instructors array, Meeting Times array, Overview, URL, Section, Subject, School, and Quarter.',
                    parameters: {
                        type: 'object',
                        properties: {
                            quarterId: {
                                type: 'integer',
                                description: 'The quarter id for which the courses will be fetched'
                            }
                        },
                        required: ['quarterId']
                    }
                }
            },
            {
                type: 'function',
                function: {
                    name: 'getCoursesBySubject',
                    description: 'Returns the list of all courses for a given subject. The courses themselves have a Title, a Number, a Topic, Overview, and Quarters array.',
                    parameters: {
                        type: 'object',
                        properties: {
                            subject: {
                                type: 'string',
                                description: 'The subject for which the courses will be fetched. Example subjects are COMP_SCI, MAT_SCI, GEN_ENG, GERMAN, ECON, ASTRON.'
                            }
                        },
                        required: ['subject']
                    }
                }
            },
            {
                type: 'function',
                function: {
                    name: 'getCoursesByKey',
                    description: 'Returns the list of all courses accross quarters for a given key. The courses themselves will have all fields. The key for a course is the course subject + number\'s first two parts.',
                    parameters: {
                        type: 'object',
                        properties: {
                            key: {
                                type: 'string',
                                description: 'The key for which the courses will be fetched. An example subject is COMP_SCI 213-0.'
                            }
                        },
                        required: ['key']
                    }
                }
            }
        ],
        tool_choice: 'auto'// the engine will decide which tool to use
    });

    while (await fulfillFunctionCalls(response, context)) {
        response = await openAI.chat.completions.create({
            model: OPENAI_MODEL,
            messages: context
        })
    }

    return (response.choices[0].message.content) || ''
}
