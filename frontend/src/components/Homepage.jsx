import React from 'react';

function Home() {
    return (
        <div className="d-flex vh-100 justify-content-center align-items-center">
            <div className="container text-center">
                <header className="my-5">
                    <div className="card custom-card p-5">
                        <h1 className="display-3 custom-title">Welcome to DocuDefense</h1>
                        <p className="lead fs-2">
                            Your Trusted Partner in Document Security and Management
                        </p>
                        <p className="py-3 fs-4">
                            DocuDefense provides secure, efficient, and reliable document management, helping you organize, protect, and access your important files whenever you need them.
                        </p>
                    </div>
                </header>

                <section id="features" className="my-5">
                    <h2 className="text-center display-4 mb-4">Why Choose DocuDefense?</h2>
                    <div className="row justify-content-center">
                        <div className="col-md-4">
                            <div className="card custom-card mb-4">
                                <div className="card-body">
                                    <h5 className="custom-title fs-4">Secure Document Storage</h5>
                                    <p className="card-text">
                                        Keep your documents safe with our advanced security protocols designed for sensitive information.
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div className="col-md-4">
                            <div className="card custom-card mb-4">
                                <div className="card-body">
                                    <h5 className="custom-title fs-4">Effortless Organization</h5>
                                    <p className="card-text">
                                        Organize files with a user-friendly tagging and folder system that makes finding documents quick and easy.
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div className="col-md-4">
                            <div className="card custom-card mb-4">
                                <div className="card-body">
                                    <h5 className="custom-title fs-4">Seamless Collaboration</h5>
                                    <p className="card-text">
                                        Share documents and assign permissions, ensuring secure and effective teamwork.
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </div>
    );
}

export default Home;
