import { createColumnHelper } from '@tanstack/table-core';
import { renderSnippet } from '$lib/components/ui/data-table';
import type { Snippet } from 'svelte';

export type SupportCase = {
	id: string;
	userName: string;
	subject: string;
	status: string;
	category: string;
	created_at: Date;
	updated_at: Date;
	resolved_at: Date | null;
};

const columnHelper = createColumnHelper<SupportCase>();

const formatDate = (date: Date) => {
	return new Date(date).toLocaleDateString('en-US', {
		month: 'short',
		day: 'numeric',
		year: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	});
};

export const createColumns = (snippets: {
	statusSnippet: Snippet<[{ status: string }]>;
}) => [
	columnHelper.accessor('subject', {
		header: 'Subject',
		cell: (info) => info.getValue()
	}),
	columnHelper.accessor('status', {
		header: 'Status',
		cell: (info) => {
			return renderSnippet(snippets.statusSnippet, { status: info.getValue() });
		}
	}),
	columnHelper.accessor('category', {
		header: 'Category',
		cell: (info) => info.getValue().replace('_', ' ')
	}),
	columnHelper.accessor('userName', {
		header: 'User',
		cell: (info) => info.getValue()
	}),
	columnHelper.accessor('created_at', {
		header: 'Created',
		cell: (info) => formatDate(info.getValue())
	}),
	columnHelper.accessor('updated_at', {
		header: 'Updated',
		cell: (info) => formatDate(info.getValue())
	})
];
