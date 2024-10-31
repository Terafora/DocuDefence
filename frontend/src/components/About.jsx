import React from 'react';

function About() {
    return (
        <div className="container my-5 text-center">
            <section className="about-section mb-5">
                <h2 className="display-5 mb-3">About Me</h2>
                <p className="fs-5">
                    Hi there! I’m Charlotte Stone, a full-stack web developer with a focus on building user-centered applications. I bring experience from technical support and software engineering, with a passion for crafting solutions that make a real impact.
                </p>
            </section>
            
            <section className="about-section mb-5">
                <h2 className="display-5 mb-3">About DocuDefense</h2>
                <p className="fs-5">
                    DocuDefense is a secure, streamlined document management platform I developed for both individual and business needs. It ensures sensitive documents are stored safely, easily accessed, and efficiently organized.
                </p>
                <p className="fs-5">
                    Perfect for legal files, personal documents, and team projects, DocuDefense provides a robust, reliable way to manage your digital files.
                </p>
            </section>
            
            <section className="about-section mb-5">
                <h2 className="display-5 mb-3">How DocuDefense Works</h2>
                <ul className="list-unstyled fs-5 text-start mx-auto" style={{ maxWidth: "600px" }}>
                    <li><strong>File Upload & Storage:</strong> Securely upload and store files for easy access.</li>
                    <li><strong>Advanced Security:</strong> Protect your documents with authentication and secure file handling.</li>
                    <li><strong>Efficient Organization:</strong> Use tagging and categorizing for fast, convenient document retrieval.</li>
                    <li><strong>User Management:</strong> Manage team access by adding or removing users with roles.</li>
                    <li><strong>Additional Features:</strong> Tools like pagination, filtering, and search enhance usability.</li>
                </ul>
                <p className="fs-5 mt-3">
                    DocuDefense is more than just storage; it’s a platform built to empower users by making document management safe, accessible, and hassle-free. I hope it’s as valuable to you as it has been in its creation!
                </p>
            </section>
        </div>
    );
}

export default About;
