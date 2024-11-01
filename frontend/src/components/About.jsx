import React from 'react';

function About() {
    return (
        <div className="container my-5 text-center">
            <div className="container">
                <div className="row g-4">
                    <section className="about-section card custom-card py-4 px-5 mb-5">
                        <h2 className="display-5 mb-3">About DocuDefense</h2>

                        <article className="fs-5 py-2">
                            <header>
                            <h3 className="pb-2">Our Purpose</h3>
                            </header>
                            <p>DocuDefense is a secure document management platform designed to protect and organize sensitive files for individuals and businesses.</p>
                        </article>
                        <hr /> 
                        <article className="fs-5 py-2">
                            <header>
                            <h3 className="pb-2">Efficiency & Productivity</h3>
                            </header>
                            <p>Built for efficiency, DocuDefense streamlines document tracking, version control, and quick access, helping you stay productive and organized.</p>
                        </article>
                        <hr />
                        <article className="fs-5 py-2">
                            <header>
                            <h3 className="pb-2">Advanced Organization</h3>
                            </header>
                            <p>With advanced search, tagging, and categorization, you can locate files instantly, no matter the volume of documents.</p>
                        </article>
                        <hr />
                        <article className="fs-5 py-2">
                            <header>
                            <h3 className="pb-2">Top-Notch Security</h3>
                            </header>
                            <p>Security is our priority. DocuDefense uses end-to-end encryption to keep your documents safe and accessible only to authorized users.</p>
                        </article>
                    </section>

                    <section className="about-section card custom-card py-4 mb-5">
                        <h2 className="display-5 mb-3">How DocuDefense Works</h2>
                        <ul className="list-unstyled fs-5 text-start mx-auto">
                            <li><strong>File Upload & Storage:</strong> Securely upload and store files for easy access.</li>
                            <li><strong>Advanced Security:</strong> Protect your documents with authentication and secure file handling.</li>
                            <li><strong>Efficient Organization:</strong> Use tagging and categorizing for fast, convenient document retrieval.</li>
                            <li><strong>User Management:</strong> Manage team access by adding or removing users with roles.</li>
                            <li><strong>Additional Features:</strong> Tools like pagination, filtering, and search enhance usability.</li>
                        </ul>
                        <hr />
                        <p className="fs-5 mt-3">
                            DocuDefense is more than just storage; it’s a platform built to empower users by making document management safe, accessible, and hassle-free. I hope it’s as valuable to you as it has been in its creation!
                        </p>
                    </section>
                </div>
            </div>
        </div>
    );
}

export default About;
