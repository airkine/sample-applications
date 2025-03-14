document.getElementById("fetchCases").addEventListener("click", async () => {
    try {
        const response = await fetch("/api/cases");
        if (!response.ok) {
            throw new Error("Failed to fetch cases");
        }
        const data = await response.json();
        const container = document.getElementById("casesContainer");
        container.innerHTML = "";

        data.cases.forEach(c => {
            const caseDiv = document.createElement("div");
            caseDiv.classList.add("case");
            caseDiv.innerHTML = `<h3>${c.title}</h3><p>${c.description}</p>`;
            container.appendChild(caseDiv);
        });
    } catch (error) {
        console.error("Error fetching cases:", error);
    }
});
