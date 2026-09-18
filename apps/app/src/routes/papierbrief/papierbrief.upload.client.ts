export const maxPaperLetterPDFBytes = 15_728_640;

export const paperLetterPDFWithinLimit = (pdf: Blob) => pdf.size <= maxPaperLetterPDFBytes;

export const putPaperLetterPDF = async (uploadUrl: string, pdf: Blob) => {
	try {
		return (
			await fetch(uploadUrl, {
				method: 'PUT',
				body: pdf,
				headers: { 'Content-Type': 'application/pdf' }
			})
		).ok;
	} catch {
		return false;
	}
};
