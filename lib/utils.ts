import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}


export function getQuarterName(tag: number): string {
    const baseTag = 4930;
    const baseYear = 2024;
    const quarters = ['Winter', 'Spring', 'Summer', 'Fall'];

    const quarterOffset = Math.floor((tag - baseTag) / 10);
    const quarterIndex = quarterOffset % 4;
    const year = baseYear + Math.floor(quarterOffset / 4);

    return `${year} ${quarters[quarterIndex]}`;
}

