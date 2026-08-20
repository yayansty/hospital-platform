import { useState } from "react";
import { Link } from "react-router-dom";
import {
    EyeIcon,
    EyeSlashIcon,
    LockClosedIcon,
    UserIcon,
} from "@heroicons/react/24/outline";

import AuthLayout from "../../layouts/AuthLayout";

export default function Login() {
    const [showPassword, setShowPassword] = useState(false);

    const [form, setForm] = useState({
        username: "",
        password: "",
        remember: false,
    });

    const handleChange = (event) => {
        const { name, value, type, checked } = event.target;

        setForm((previous) => ({
            ...previous,
            [name]: type === "checkbox" ? checked : value,
        }));
    };

    const handleSubmit = (event) => {
        event.preventDefault();

        console.log("Login form:", form);
    };

    return (
        <AuthLayout>

            {/* Mobile Brand */}
            <div className="mb-8 flex items-center gap-3 lg:hidden">
                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-primary-600 text-white shadow-sm">
                    <span className="text-xl font-bold">
                        H
                    </span>
                </div>

                <div>
                    <h1 className="text-base font-semibold text-slate-900">
                        Hospital Information System
                    </h1>

                    <p className="text-xs text-slate-500">
                        Integrated Healthcare Platform
                    </p>
                </div>
            </div>

            {/* Heading */}
            <div className="mb-8">
                <h2 className="text-2xl font-bold tracking-tight text-slate-900 sm:text-3xl">
                    Welcome back
                </h2>

                <p className="mt-2 text-sm leading-6 text-slate-500">
                    Silakan masuk menggunakan akun Anda untuk melanjutkan.
                </p>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit} className="space-y-5">

                {/* Username */}
                <div>
                    <label
                        htmlFor="username"
                        className="mb-2 block text-sm font-medium text-slate-700"
                    >
                        Username
                    </label>

                    <div className="relative">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                            <UserIcon className="h-5 w-5 text-slate-400" />
                        </div>

                        <input
                            id="username"
                            name="username"
                            type="text"
                            autoComplete="username"
                            value={form.username}
                            onChange={handleChange}
                            placeholder="Masukkan username"
                            className="block w-full rounded-xl border border-slate-300 bg-white py-3 pl-11 pr-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-primary-500 focus:ring-4 focus:ring-primary-100"
                        />
                    </div>
                </div>

                {/* Password */}
                <div>
                    <label
                        htmlFor="password"
                        className="mb-2 block text-sm font-medium text-slate-700"
                    >
                        Password
                    </label>

                    <div className="relative">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
                            <LockClosedIcon className="h-5 w-5 text-slate-400" />
                        </div>

                        <input
                            id="password"
                            name="password"
                            type={showPassword ? "text" : "password"}
                            autoComplete="current-password"
                            value={form.password}
                            onChange={handleChange}
                            placeholder="Masukkan password"
                            className="block w-full rounded-xl border border-slate-300 bg-white py-3 pl-11 pr-12 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-primary-500 focus:ring-4 focus:ring-primary-100"
                        />

                        <button
                            type="button"
                            onClick={() =>
                                setShowPassword((previous) => !previous)
                            }
                            className="absolute inset-y-0 right-0 flex items-center pr-3.5 text-slate-400 transition hover:text-slate-600"
                            aria-label={
                                showPassword
                                    ? "Sembunyikan password"
                                    : "Tampilkan password"
                            }
                        >
                            {showPassword ? (
                                <EyeSlashIcon className="h-5 w-5" />
                            ) : (
                                <EyeIcon className="h-5 w-5" />
                            )}
                        </button>
                    </div>
                </div>

                {/* Remember + Forgot */}
                <div className="flex items-center justify-between gap-4">

                    <label className="flex cursor-pointer items-center gap-2">
                        <input
                            type="checkbox"
                            name="remember"
                            checked={form.remember}
                            onChange={handleChange}
                            className="h-4 w-4 rounded border-slate-300 text-primary-600 focus:ring-primary-500"
                        />

                        <span className="text-sm text-slate-600">
                            Remember me
                        </span>
                    </label>

                    <Link
                        to="/forgot-password"
                        className="text-sm font-medium text-primary-600 transition hover:text-primary-700"
                    >
                        Lupa password?
                    </Link>

                </div>

                {/* Login Button */}
                <button
                    type="submit"
                    className="flex w-full items-center justify-center rounded-xl bg-primary-600 px-4 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-primary-700 focus:outline-none focus:ring-4 focus:ring-primary-100 active:scale-[0.99]"
                >
                    Login
                </button>

            </form>

            {/* Footer */}
            <div className="mt-8 text-center">
                <p className="text-xs leading-5 text-slate-400">
                    Sistem informasi terintegrasi untuk mendukung pelayanan
                    kesehatan.
                </p>
            </div>

        </AuthLayout>
    );
}