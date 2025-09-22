// src/components/ui/table.tsx
import * as React from "react";
import { cn } from "@/lib/utils"; // your utility for classNames

export const Table = ({ className, ...props }: React.TableHTMLAttributes<HTMLTableElement>) => (
  <div className="relative w-full overflow-x-auto">
    <table className={cn("w-full text-sm border-collapse", className)} {...props} />
  </div>
);

export const TableHeader = ({ children, className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) => (
  <thead className={cn("bg-gray-100", className)} {...props}>
    {children}
  </thead>
);

export const TableBody = ({ children, className, ...props }: React.HTMLAttributes<HTMLTableSectionElement>) => (
  <tbody className={cn("", className)} {...props}>
    {children}
  </tbody>
);

export const TableRow = ({ children, className, ...props }: React.HTMLAttributes<HTMLTableRowElement>) => (
  <tr className={cn("hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors", className)} {...props}>
    {children}
  </tr>
);

export const TableHeadCell = ({ children, className, ...props }: React.ThHTMLAttributes<HTMLTableCellElement>) => (
  <th className={cn("px-4 py-2 text-left font-medium text-gray-700 dark:text-gray-200", className)} {...props}>
    {children}
  </th>
);

export const TableCell = ({ children, className, ...props }: React.TdHTMLAttributes<HTMLTableCellElement>) => (
  <td className={cn("px-4 py-2 border-b border-gray-200 dark:border-gray-700", className)} {...props}>
    {children}
  </td>
);
