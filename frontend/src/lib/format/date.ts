import { format, parseISO, differenceInCalendarDays } from 'date-fns';

export function formatDate(iso: string, pattern = 'MMM d, yyyy'): string {
  return format(parseISO(iso), pattern);
}

export function formatDateRange(startIso: string, endIso: string): string {
  return `${formatDate(startIso)} – ${formatDate(endIso)}`;
}

export function daysBetween(startIso: string, endIso: string): number {
  return differenceInCalendarDays(parseISO(endIso), parseISO(startIso));
}
