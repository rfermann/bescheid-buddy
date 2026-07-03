<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import { sendOTP, verifyOTP } from '../auth.remote';

	let step = $state(1);

	sendOTP.enhance(async ({ submit }) => {
		await submit();
	});

	verifyOTP.enhance(async ({ submit }) => {
		await submit();
	});

	$effect(() => {
		if (sendOTP.result?.success) {
			step = 2;
		}
	});

	let submitIssues = $derived.by(() => {
		return verifyOTP.fields.allIssues()?.filter((issue) => issue.path.length === 0);
	});

</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50">
	<Card class="w-full max-w-md">
		<CardHeader>
			<CardTitle>Admin Login</CardTitle>
			<CardDescription>
				{step === 1 ? 'Enter your email to receive a one-time password for secure login.' : 'Check your email for the verification code.'}
			</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			{#if sendOTP.result?.message}
				<div
					class="text-sm p-3 rounded-md {sendOTP.result.success ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50'}">
					{sendOTP.result.message}
				</div>
			{/if}

			<div class="space-y-2">
				<Label for="email">Email</Label>
				{#if step === 1}
					<form {...sendOTP}>
						<Input
							id="email"
							placeholder="admin@example.com"
							{...sendOTP.fields.email.as('email')}
						/>
						<Button
							type="submit"
							disabled={!!sendOTP.pending}
							class="w-full mt-4"
						>
							{sendOTP.pending ? 'Sending...' : 'Send OTP'}
						</Button>
					</form>
				{:else}
					<Input
						id="email"
						type="email"
						value={sendOTP.fields.email.value()}
						disabled
					/>
					<form {...sendOTP}>
						<input {...sendOTP.fields.email.as('hidden', sendOTP.fields.email.value() || '')} />
						<Button
							type="submit"
							variant="outline"
							disabled={!!sendOTP.pending}
							class="w-full mt-2"
						>
							{sendOTP.pending ? 'Resending...' : 'Resend Code'}
						</Button>
					</form>
				{/if}
			</div>

			{#if step === 2}
				<form {...verifyOTP} class="space-y-4">
					{#if submitIssues}
						{#each submitIssues as issue}
							<div class="text-sm text-red-600 bg-red-50 p-3 rounded-md">
								{issue.message}
							</div>
						{/each}
					{/if}

					{#if verifyOTP.result?.message}
						<div
							class="text-sm p-3 rounded-md {verifyOTP.result.success ? 'text-green-600 bg-green-50' : 'text-red-600 bg-red-50'}">
							{verifyOTP.result.message}
						</div>
					{/if}

					<input {...verifyOTP.fields.email.as('hidden', sendOTP.fields.email.value())} />

					<div class="space-y-2">
						<Label for="otp">One-Time Password</Label>
						<Input
							id="otp"
							type="text"
							placeholder="Enter 6-digit code"
							{...verifyOTP.fields.otp.as('text')}
							maxlength="6"
							inputmode="numeric"
							pattern="[0-9]*"
						/>
					</div>

					<Button
						type="submit"
						disabled={!!verifyOTP.pending}
						class="w-full"
					>
						{verifyOTP.pending ? 'Verifying...' : 'Verify OTP'}
					</Button>

					<Button
						type="button"
						variant="outline"
						onClick={() => step = 1}
						class="w-full"
					>
						Back to Email
					</Button>
				</form>
			{/if}
		</CardContent>
	</Card>
</div>
