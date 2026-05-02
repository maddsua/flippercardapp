<script setup lang="ts">
import Card from './components/Cards/Card.vue';

import type { ContentNode } from './components/Cards/content';


const content: ContentNode[] = [
	{
		type: 'title',
		content: 'Card title'
	},
	{
		type: 'textbox',
		content: [
			{
				type: 'text',
				content: 'Sample text to ask'
			},
			{
				type: 'text',
				content: ' a question?'
			},
			{
				type: 'newline'
			},
			{
				type: 'text',
				content: 'Your answer?'
			}
		]
	},
	{
		type: 'poll',
		content: [
			{
				value: 'Option 1',
				action: 'fail-show',
			},
			{
				value: 'Option 2',
				action: 'fail-show',
			},
			{
				value: 'Option 2',
				action: 'next',
			}
		]
	}
];

</script>

<template>
	<Card :front="{ content }" :back="{ content: []}" />
</template>
