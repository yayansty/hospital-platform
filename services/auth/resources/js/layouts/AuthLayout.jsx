export default function AuthLayout({ children }) {
    return (
        <div className="min-h-screen bg-slate-100">
            <div className="flex min-h-screen">

                {/* Left Panel */}
                <div className="relative hidden overflow-hidden bg-primary-700 lg:flex lg:w-1/2">
                    <div className="absolute inset-0 bg-gradient-to-br from-primary-700 via-primary-800 to-primary-900" />

                    <div className="relative z-10 flex w-full flex-col justify-between p-10 xl:p-16">

                        {/* Brand */}
                        <div>
                            <div className="flex items-center gap-3">
                                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-white/10 text-white ring-1 ring-white/20">
                                    <span className="text-xl font-bold">
                                        H
                                    </span>
                                </div>

                                <div>
                                    <h1 className="text-lg font-semibold text-white">
                                        Hospital Information System
                                    </h1>

                                    <p className="text-sm text-primary-100">
                                        Integrated Healthcare Platform
                                    </p>
                                </div>
                            </div>
                        </div>

                        {/* Information */}
                        <div className="max-w-xl">
                            <h2 className="text-4xl font-bold leading-tight text-white xl:text-5xl">
                                Manage healthcare
                                <br />
                                smarter and better.
                            </h2>

                            <p className="mt-6 max-w-lg text-base leading-7 text-primary-100 xl:text-lg">
                                Satu platform terintegrasi untuk mendukung
                                operasional dan pelayanan rumah sakit secara
                                efektif, aman, dan terukur.
                            </p>
                        </div>

                        {/* Footer */}
                        <div className="text-sm text-primary-200">
                            © {new Date().getFullYear()} Hospital Information
                            System
                        </div>
                    </div>
                </div>

                {/* Right Panel */}
                <div className="flex w-full items-center justify-center px-5 py-10 sm:px-8 lg:w-1/2 lg:px-12 xl:px-20">
                    <div className="w-full max-w-md">
                        {children}
                    </div>
                </div>

            </div>
        </div>
    );
}