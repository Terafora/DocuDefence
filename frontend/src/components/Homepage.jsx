import React from 'react';

function Home() {
    return (
        <div>
            <header className="container text-center my-5">
                <div className="card custom-card shadow-lg p-5">
                    <h1 className="display-4 custom-title">Welcome to DocuDefense</h1>
                    <p className="lead">
                        Your Trusted Partner in Document Security and Management
                    </p>
                    <p>
                        DocuDefense provides secure, efficient, and reliable document management, helping you organize, protect, and access your important files whenever you need them.
                    </p>
                </div>
            </header>

            <section id="features" className="container my-5">
                <h2 className="text-center custom-title mb-4">Why Choose DocuDefense?</h2>
                <div className="row">
                    <div className="col-md-4">
                        <div className="card custom-card mb-4 shadow-sm">
                            <div className="card-body">
                                <h5 className="custom-title">Secure Document Storage</h5>
                                <p className="card-text">Keep your documents safe with our advanced security protocols designed for sensitive information.</p>
                            </div>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card custom-card mb-4 shadow-sm">
                            <div className="card-body">
                                <h5 className="custom-title">Effortless Organization</h5>
                                <p className="card-text">Organize files with a user-friendly tagging and folder system that makes finding documents quick and easy.</p>
                            </div>
                        </div>
                    </div>
                    <div className="col-md-4">
                        <div className="card custom-card mb-4 shadow-sm">
                            <div className="card-body">
                                <h5 className="custom-title">Seamless Collaboration</h5>
                                <p className="card-text">Share documents and assign permissions, ensuring secure and effective teamwork.</p>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <footer id="contact" className="card custom-card text-center py-4">
                <div className="container">
                    <p className="mb-0">Need help or have questions? <a href="mailto:support@docudefense.com">Contact Support</a></p>
                    <p>&copy; 2023 DocuDefense. All rights reserved.</p>
                </div>
            </footer>
        </div>
    );
}

export default Home;
