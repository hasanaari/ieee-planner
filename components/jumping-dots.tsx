"use client"
import { cn } from "@/lib/utils"

interface JumpingDotsProps {
    className?: string
    dotClassName?: string
    dotSize?: number
    dotCount?: number
    speed?: "slow" | "normal" | "fast"
    color?: string
}

export function JumpingDots({
    className,
    dotClassName,
    dotSize = 8,
    dotCount = 3,
    speed = "normal",
    color = "currentColor",
}: JumpingDotsProps) {
    // Calculate animation duration based on speed
    const getDuration = () => {
        switch (speed) {
            case "slow":
                return 1.4
            case "fast":
                return 0.7
            default:
                return 1
        }
    }

    const duration = getDuration()

    return (
        <div className={cn("flex items-center justify-center space-x-2", className)}>
            {Array.from({ length: dotCount }).map((_, index) => (
                <div
                    key={index}
                    className={cn("rounded-full animate-jump", dotClassName)}
                    style={{
                        width: `${dotSize}px`,
                        height: `${dotSize}px`,
                        backgroundColor: color,
                        animation: `jump ${duration}s infinite`,
                        animationDelay: `${index * (duration / 10)}s`,
                    }}
                />
            ))}

            <style jsx global>{`
        @keyframes jump {
          0%, 100% {
            transform: translateY(0);
          }
          50% {
            transform: translateY(-10px);
          }
        }
        
        .animate-jump {
          animation: jump 1s infinite;
        }
      `}</style>
        </div>
    )
}

