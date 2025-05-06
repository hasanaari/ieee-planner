"use client"

import { useEffect, useRef, useState } from "react"
import { Send } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { JumpingDots } from "@/components/jumping-dots"
import { askChat, type Message } from "../app/actions/chat"

// TODO: I don't know chat officially uses as its response syntax but its response in raw text are rendered funny, I think we are better off with an md renderer
export function ChatSection({ selectedCourses, selectedMajor, allquarters, messages, setMessages }: { selectedCourses: string[], selectedMajor: string, allquarters: number[], messages: Message[], setMessages: React.Dispatch<React.SetStateAction<Message[]>> }) {
    const bottomRef = useRef<HTMLDivElement>(null);
    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [messages]);
    const [inputValue, setInputValue] = useState("")
    const [inputLocked, setInputLocked] = useState(false)

    const handleSendMessage = async () => {
        if (!inputValue.trim()) return

        const userMessage: Message = {
            id: Date.now().toString(),
            content: inputValue,
            sender: "user",
            timestamp: new Date(),
            isError: false
        }

        setMessages(prev => [...prev, userMessage]);
        setInputValue("");

        setInputLocked(true);
        try {
            const resp = await askChat(selectedCourses, selectedMajor, allquarters, inputValue);
            if (resp) {
                const msg: Message = {
                    id: Date.now().toString(),
                    content: resp,
                    sender: "assistant",
                    timestamp: new Date(),
                    isError: false
                }
                setMessages(prev => [...prev, msg])
            } else {
                const noMsgErr: Message = {
                    id: Date.now().toString(),
                    content: "There has been an unexpected error in getting a response back!",
                    sender: "assistant",
                    timestamp: new Date(),
                    isError: true
                }
                setMessages(prev => [...prev, noMsgErr])
            }
        } catch (error) {
            console.error("Chat error:", error);
            const err: Message = {
                id: Date.now().toString(),
                content: "There has been an unexpected error: " + JSON.stringify(error),
                sender: "assistant",
                timestamp: new Date(),
                isError: true
            }

            setMessages(prev => [...prev, err])
            setInputLocked(false);
        } finally {
            setInputLocked(false);
        }
    }

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault()
            handleSendMessage()
        }
    }

    return (
        <div className="flex flex-col h-[calc(100vh-12rem)]">
            <h2 className="text-lg font-semibold mb-4">Chat Assistant</h2>

            <ScrollArea className="flex-1 p-4 border rounded-md mb-4">
                <div className="space-y-4">
                    {messages.map((message) => (
                        <div
                            key={message.id}
                            className={`flex ${message.sender === "user" ? "justify-end" : "justify-start"}`}
                        >
                            <div
                                className={`max-w-[80%] rounded-lg px-4 py-2 ${message.isError
                                    ? "bg-destructive text-destructive-foreground"
                                    : message.sender === "user"
                                        ? "bg-primary text-primary-foreground"
                                        : "bg-muted"
                                    }`}
                            >
                                <p>{message.content}</p>
                                <p className="text-xs opacity-70 mt-1">
                                    {message.timestamp.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
                {inputLocked &&
                    <div className="bg-gray-100 dark:bg-gray-700 p-3 rounded-lg max-w-[5%]">
                        <JumpingDots dotSize={6} />
                    </div>
                }
                <div ref={bottomRef} />
            </ScrollArea>

            <div className="flex items-center space-x-2">
                <Input
                    placeholder="Type your message..."
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyDown={handleKeyDown}
                    className="flex-1"
                    disabled={inputLocked}
                />
                <Button onClick={handleSendMessage} size="icon">
                    <Send className="h-4 w-4" />
                    <span className="sr-only">Send message</span>
                </Button>
            </div>
        </div>
    )
}

