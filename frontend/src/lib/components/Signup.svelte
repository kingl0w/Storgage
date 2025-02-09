<script lang="ts">
    import { goto } from '$app/navigation';
    import { signup } from '$lib/api';

    let username = '';
    let password = '';
    let confirmPassword = '';
    let inviteCode = '';
    let error = '';

    async function handleSignup() {
        if (password !== confirmPassword) {
            error = 'Passwords do not match';
            return;
        }

        
        const response = await signup(username, password, inviteCode);

        if (response.ok) {
            await goto('/login');
        } else {
            let data;
            try {
                data = await response.json();
            } catch (e) {
                const textErr = await response.text();
                console.error('Server error:', textErr);
                error = textErr || 'Signup failed';
                return;
            }
            error = data.error || 'Signup failed';
            console.error('Signup error:', data);
        }
    }
</script>

<div class="signup-container">
    <h1>Sign Up</h1>
    
    {#if error}
        <div class="error">{error}</div>
    {/if}
    
    <form on:submit|preventDefault={handleSignup}>
        <div class="input-group">
            <label for="username">Username</label>
            <input 
                type="text" 
                id="username" 
                bind:value={username} 
                required
            />
        </div>
        
        <div class="input-group">
            <label for="password">Password</label>
            <input 
                type="password" 
                id="password" 
                bind:value={password} 
                required
            />
        </div>

        <div class="input-group">
            <label for="confirmPassword">Confirm Password</label>
            <input 
                type="password" 
                id="confirmPassword" 
                bind:value={confirmPassword} 
                required
            />
        </div>

        <div class="input-group">
            <label for="inviteCode">Invitation Code</label>
            <input 
                type="text" 
                id="inviteCode" 
                bind:value={inviteCode} 
                required
            />
        </div>
        
        <button type="submit">Sign Up</button>
    </form>

    <div class="login-link">
        Already have an account? <a href="/login">Login here</a>
    </div>
</div>

<style>
    .signup-container {
        max-width: 400px;
        margin: 100px auto;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        background: white;
    }

    .input-group {
        margin-bottom: 1rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
    }

    input {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #ccc;
        border-radius: 4px;
    }

    button {
        width: 100%;
        padding: 0.75rem;
        background: #4a9eff;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    button:hover {
        background: #2b7fd9;
    }

    .error {
        color: red;
        margin-bottom: 1rem;
        padding: 0.5rem;
        background: #fee;
        border-radius: 4px;
    }

    .login-link {
        margin-top: 1rem;
        text-align: center;
        font-size: 0.9rem;
    }

    .login-link a {
        color: #4a9eff;
        text-decoration: none;
    }

    .login-link a:hover {
        text-decoration: underline;
    }
</style>
