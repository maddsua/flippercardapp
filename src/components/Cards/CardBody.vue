<template>
	<div class="card-body">
		<div class="card-content">
			<slot>
				[Card content]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.card-body {
		position: absolute;
		top: 0;
		bottom: 0;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		width: 100%;
		height: 100%;
		background-color: white;
		border-radius: 2em;
		padding: 1em;
		
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;

		// these two somehow prevent jagged edges
		outline: 1px solid transparent;
		box-shadow: 1em 0.25rem 0.25rem rgba(0, 0, 0, 0.25);

		&:nth-child(2n) {
			transform: rotateY(180deg);
		}

		&:not(:nth-child(2n))::before {
			content: "?";
			display: block;
			font-size: 3.5em;
			line-height: 1em;
			position: absolute;
			left: 0.65em;
			top: 0.5em;
			font-weight: 600;
			color: orange;
		}
	}
	.card-content {
		display: flex;
		flex-direction: column;
		gap: 2em;
		align-items: center;
		flex-grow: 1;
		border: 3px solid orange;
		border-radius: 1.5em;
		padding: 2em;
		color: black;
	}
</style>
