import React, { useEffect, useRef } from 'react';
import { getDocument, GlobalWorkerOptions } from 'pdfjs-dist';

// Set workerSrc to correctly point to the worker in the public folder
GlobalWorkerOptions.workerSrc = `${process.env.PUBLIC_URL}/pdf.worker.min.js`;

const PDFPreview = ({ fileBlob }) => {
    const canvasRef = useRef(null);

    useEffect(() => {
        const loadPdf = async () => {
            if (fileBlob) {
                try {
                    // Ensure the fileBlob is valid binary data
                    const pdf = await getDocument({ data: fileBlob }).promise;
                    const page = await pdf.getPage(1);
                    const viewport = page.getViewport({ scale: 1.5 });

                    const canvas = canvasRef.current;
                    const context = canvas.getContext('2d');
                    canvas.width = viewport.width;
                    canvas.height = viewport.height;

                    const renderContext = {
                        canvasContext: context,
                        viewport: viewport,
                    };
                    await page.render(renderContext).promise;
                } catch (error) {
                    console.error("Error rendering PDF:", error);
                }
            }
        };

        loadPdf();
    }, [fileBlob]);

    return <canvas ref={canvasRef} />;
};

export default PDFPreview;
